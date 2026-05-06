package shared

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

// based on https://github.com/refraction-networking/utls/issues/16#issuecomment-1285198375
func NewBypassJA3Transport(helloID utls.ClientHelloID, disableHTTP2 bool) *BypassJA3Transport {
	return &BypassJA3Transport{clientHello: helloID, disableHTTP2: disableHTTP2}
}

type BypassJA3Transport struct {
	tr1 http.Transport
	tr2 http2.Transport

	mu           sync.RWMutex
	clientHello  utls.ClientHelloID
	disableHTTP2 bool
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

	conn, err := net.Dial("tcp", net.JoinHostPort(req.URL.Host, port))
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
	nextProtos := []string{}
	if !b.disableHTTP2 {
		nextProtos = []string{"h2"}
	}
	return &utls.Config{
		ServerName:         req.URL.Host,
		InsecureSkipVerify: true,
		NextProtos:         nextProtos,
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

// BrowserHeaderRoundTripper injects browser-like headers into all requests
type BrowserHeaderRoundTripper struct {
	transport http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface by injecting browser headers
func (b *BrowserHeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Inject headers only if not already present (allows per-request override)
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", UserAgent)
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	}
	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "en-GB,en;q=0.9")
	}
	// if req.Header.Get("Accept-Encoding") == "" {
	// 	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	// }
	if req.Header.Get("Sec-Fetch-Mode") == "" {
		req.Header.Set("Sec-Fetch-Mode", "navigate")
	}
	if req.Header.Get("Sec-Fetch-Site") == "" {
		req.Header.Set("Sec-Fetch-Site", "none")
	}
	if req.Header.Get("Sec-Fetch-Dest") == "" {
		req.Header.Set("Sec-Fetch-Dest", "document")
	}
	if req.Header.Get("sec-ch-ua") == "" {
		req.Header.Set("sec-ch-ua", `"Google Chrome";v="147", "Not.A/Brand";v="8", "Chromium";v="147"`)
	}

	return b.transport.RoundTrip(req)
}

func NewClient(disableHTTP2 bool) *http.Client {
	jar, _ := cookiejar.New(nil)
	return NewClientWithJar(jar, disableHTTP2)
}

func NewClientWithJar(jar *cookiejar.Jar, disableHTTP2 bool) *http.Client {
	baseTransport := NewBypassJA3Transport(utls.HelloChrome_Auto, disableHTTP2)
	return &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
		Transport: &BrowserHeaderRoundTripper{
			transport: baseTransport,
		},
	}
}
