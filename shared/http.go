package shared

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"slices"
	"time"

	utls "github.com/refraction-networking/utls"
)

type utlsConn struct {
	*utls.UConn
}

func (c *utlsConn) ConnectionState() tls.ConnectionState {
	uState := c.UConn.ConnectionState()
	return tls.ConnectionState{
		Version:                     uState.Version,
		HandshakeComplete:           uState.HandshakeComplete,
		DidResume:                   uState.DidResume,
		CipherSuite:                 uState.CipherSuite,
		NegotiatedProtocol:          uState.NegotiatedProtocol,
		NegotiatedProtocolIsMutual:  uState.NegotiatedProtocolIsMutual,
		ServerName:                  uState.ServerName,
		PeerCertificates:            uState.PeerCertificates,
		VerifiedChains:              uState.VerifiedChains,
		SignedCertificateTimestamps: uState.SignedCertificateTimestamps,
		OCSPResponse:                uState.OCSPResponse,
		TLSUnique:                   uState.TLSUnique,
	}
}

// BrowserHeaderRoundTripper injects browser-like headers into all requests
type BrowserHeaderRoundTripper struct {
	transport http.RoundTripper
}

func NewBrowserHeaderRoundTripper(transport http.RoundTripper) *BrowserHeaderRoundTripper {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &BrowserHeaderRoundTripper{transport: transport}
}

// RoundTrip implements the http.RoundTripper interface by injecting browser headers
func (b *BrowserHeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
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
		req.Header.Set("sec-ch-ua", `"Google Chrome";v="153", "Not.A/Brand";v="8", "Chromium";v="153"`)
	}
	if req.Header.Get("sec-ch-ua-mobile") == "" {
		req.Header.Set("sec-ch-ua-mobile", "?0")
	}
	if req.Header.Get("sec-ch-ua-platform") == "" {
		req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	}
	if req.Header.Get("Upgrade-Insecure-Requests") == "" {
		req.Header.Set("Upgrade-Insecure-Requests", "1")
	}

	tr := b.transport
	if tr == nil {
		tr = http.DefaultTransport
	}
	return tr.RoundTrip(req)
}

func NewClient(nextProtos []string) *http.Client {
	jar, _ := cookiejar.New(nil)
	if len(nextProtos) == 0 {
		return NewClientWithJar(jar, []string{"h2", "http/1.1"})
	}
	return NewClientWithJar(jar, nextProtos)
}

func NewClientWithJar(jar *cookiejar.Jar, nextProtos []string) *http.Client {
	if len(nextProtos) == 0 {
		nextProtos = []string{"h2", "http/1.1"}
	}

	hasHTTP2 := slices.Contains(nextProtos, "h2")
	hasHTTP1 := slices.Contains(nextProtos, "http/1.1")

	protocols := new(http.Protocols)
	protocols.SetHTTP1(hasHTTP1)
	protocols.SetHTTP2(hasHTTP2)

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		Protocols:             protocols,
		ForceAttemptHTTP2:     hasHTTP2,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				host = addr
			}

			rawConn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}

			config := &utls.Config{
				ServerName: host,
				NextProtos: nextProtos,
			}

			spec, err := utls.UTLSIdToSpec(utls.HelloChrome_Auto)
			if err != nil {
				_ = rawConn.Close()
				return nil, err
			}

			for _, ext := range spec.Extensions {
				if alpn, ok := ext.(*utls.ALPNExtension); ok {
					alpn.AlpnProtocols = nextProtos
				}
			}

			uConn := utls.UClient(rawConn, config, utls.HelloCustom)
			if err := uConn.ApplyPreset(&spec); err != nil {
				_ = rawConn.Close()
				return nil, err
			}

			if err := uConn.HandshakeContext(ctx); err != nil {
				_ = rawConn.Close()
				return nil, err
			}

			return &utlsConn{UConn: uConn}, nil
		},
	}

	return &http.Client{
		Jar:       jar,
		Timeout:   30 * time.Second,
		Transport: NewBrowserHeaderRoundTripper(transport),
	}
}
