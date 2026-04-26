package shared

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"sync"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

// based on https://github.com/refraction-networking/utls/issues/16#issuecomment-1285198375
func NewBypassJA3Transport(helloID utls.ClientHelloID) *BypassJA3Transport {
	return &BypassJA3Transport{clientHello: helloID}
}

type BypassJA3Transport struct {
	tr1 http.Transport
	tr2 http2.Transport

	mu          sync.RWMutex
	clientHello utls.ClientHelloID
}

type responseBodyCloser struct {
	io.ReadCloser
	closeFn func() error
	once    sync.Once
}

func (b *responseBodyCloser) Close() error {
	err := b.ReadCloser.Close()
	b.once.Do(func() {
		if closeErr := b.closeFn(); err == nil && closeErr != nil {
			err = closeErr
		}
	})
	return err
}

func (b *BypassJA3Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch req.URL.Scheme {
	case "https":
		return b.httpsRoundTrip(req)
	case "http":
		return b.tr1.RoundTrip(req)
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", req.URL.Scheme)
	}
}

func (b *BypassJA3Transport) httpsRoundTrip(req *http.Request) (*http.Response, error) {
	port := req.URL.Port()
	if port == "" {
		port = "443"
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", req.URL.Host, port))
	if err != nil {
		return nil, fmt.Errorf("tcp net dial fail: %w", err)
	}

	tlsConn, err := b.tlsConnect(conn, req)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("tls connect fail: %w", err)
	}

	httpVersion := tlsConn.ConnectionState().NegotiatedProtocol
	switch httpVersion {
	case "h2":
		clientConn, err := b.tr2.NewClientConn(tlsConn)
		if err != nil {
			_ = tlsConn.Close()
			return nil, fmt.Errorf("create http2 client with connection fail: %w", err)
		}

		resp, err := clientConn.RoundTrip(req)
		if err != nil {
			_ = clientConn.Close()
			return nil, err
		}
		resp.Body = &responseBodyCloser{ReadCloser: resp.Body, closeFn: clientConn.Close}
		return resp, nil
	case "http/1.1", "":
		err := req.Write(tlsConn)
		if err != nil {
			_ = tlsConn.Close()
			return nil, fmt.Errorf("write http1 tls connection fail: %w", err)
		}

		resp, err := http.ReadResponse(bufio.NewReader(tlsConn), req)
		if err != nil {
			_ = tlsConn.Close()
			return nil, err
		}
		resp.Body = &responseBodyCloser{ReadCloser: resp.Body, closeFn: tlsConn.Close}
		return resp, nil
	default:
		_ = tlsConn.Close()
		return nil, fmt.Errorf("unsuported http version: %s", httpVersion)
	}
}

func (b *BypassJA3Transport) getTLSConfig(req *http.Request) *utls.Config {
	return &utls.Config{
		ServerName:         req.URL.Host,
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2"},
	}
}

func (b *BypassJA3Transport) tlsConnect(conn net.Conn, req *http.Request) (*utls.UConn, error) {
	b.mu.RLock()
	tlsConn := utls.UClient(conn, b.getTLSConfig(req), b.clientHello)
	b.mu.RUnlock()

	if err := tlsConn.Handshake(); err != nil {
		return nil, fmt.Errorf("tls handshake fail: %w", err)
	}
	return tlsConn, nil
}

func (b *BypassJA3Transport) SetClientHello(hello utls.ClientHelloID) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clientHello = hello
}

type VSCOPost struct {
	Medias struct {
		ByID map[string]struct {
			Media struct {
				PermaSubdomain string `json:"permaSubdomain"`
				ResponsiveURL  string `json:"responsiveUrl"`
				VideoURL       string `json:"videoUrl"`
				PlaybackURL    string `json:"playbackUrl"`
				PosterURL      string `json:"posterUrl"`
				PosterWidth    uint   `json:"widthPx"`
				Site           struct {
					Domain string `json:"domain"`
				} `json:"site"`
			} `json:"media"`
		} `json:"byId"`
	} `json:"medias"`
}

func findFirstURL(response io.ReadCloser) string {
	scanner := bufio.NewScanner(response)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "https://") {
			return line
		}
	}
	return ""
}

func extractStreamURL(playbackURL string) (string, error) {
	playlistResponse, err := http.Get(playbackURL)
	if err != nil {
		return "", err
	}
	defer playlistResponse.Body.Close()
	renditionURL := findFirstURL(playlistResponse.Body)
	if len(renditionURL) == 0 {
		return "", errors.New("couldn't find rendition URL")
	}
	renditionResponse, err := http.Get(renditionURL)
	if err != nil {
		return "", err
	}
	defer renditionResponse.Body.Close()
	streamURL := findFirstURL(renditionResponse.Body)
	if len(streamURL) == 0 {
		return "", errors.New("couldn't find stream URL")
	}
	return streamURL, nil
}

var vsco_regexp = regexp.MustCompile(`<script>window\.__PRELOADED_STATE__ =(.*?)</script>`)

func VSCO(owner, post string) ([]string, string, []*http.Cookie, error) {
	postURL := fmt.Sprintf("https://vsco.co/%s/media/%s", owner, post)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}

	client := &http.Client{
		Jar:       jar,
		Timeout:   time.Second * 30,
		Transport: NewBypassJA3Transport(utls.HelloChrome_Auto),
	}

	htmlRequest, err := http.NewRequest(http.MethodGet, postURL, nil)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}
	htmlRequest.Header.Add("User-Agent", UserAgent)
	htmlRequest.Header.Add("sec-ch-ua", `"Google Chrome";v="147", "Not.A/Brand";v="8", "Chromium";v="147"`)
	htmlRequest.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;image/avif,image/webp,image/apng,*/*;application/signed-exchange;")

	htmlResponse, err := client.Do(htmlRequest)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}
	defer htmlResponse.Body.Close()

	body, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}

	script := vsco_regexp.FindString(string(body))
	if script == "" {
		return []string{}, "", []*http.Cookie{}, errors.New("could not find JSON")
	}

	jsonText := script[len("<script>window.__PRELOADED_STATE__ =") : len(script)-len("</script>")]
	jsonText = strings.ReplaceAll(jsonText, "undefined", "null")
	var vscoPost VSCOPost
	if err := json.Unmarshal([]byte(jsonText), &vscoPost); err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}

	media := vscoPost.Medias.ByID[post]
	username := media.Media.PermaSubdomain
	URLs := make([]string, 0, 2)

	if len(media.Media.VideoURL) > 0 {
		URLs = append(URLs, fmt.Sprintf("https://%s", media.Media.VideoURL))
	} else if len(media.Media.PlaybackURL) > 0 {
		username = media.Media.Site.Domain
		URL, err := extractStreamURL(media.Media.PlaybackURL)
		if err != nil {
			return []string{}, "", []*http.Cookie{}, err
		}
		URLs = append(URLs, URL)
		URLs = append(URLs, fmt.Sprintf("%s?time=0&width=%d", media.Media.PosterURL, media.Media.PosterWidth))
	} else {
		URL := fmt.Sprintf("https://%s", media.Media.ResponsiveURL)
		URLs = append(URLs, URL)
	}

	return URLs, username, jar.Cookies(htmlResponse.Request.URL), err
}
