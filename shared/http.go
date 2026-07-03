package shared

import (
	"context"
	"crypto/tls"
	"fmt"
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

func (b *BypassJA3Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch req.URL.Scheme {
	case "https":
		if b.disableHTTP2 {
			return b.tr1.RoundTrip(req)
		}
		return b.tr2.RoundTrip(req)
	case "http":
		return b.tr1.RoundTrip(req)
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", req.URL.Scheme)
	}
}

func (b *BypassJA3Transport) getTLSConfig(serverName string, nextProtos []string) *utls.Config {
	return &utls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: true,
		NextProtos:         nextProtos,
	}
}

func (b *BypassJA3Transport) tlsConnect(ctx context.Context, conn net.Conn, serverName string, nextProtos []string) (*utls.UConn, error) {
	// Set deadline on the underlying connection based on context
	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			return nil, fmt.Errorf("failed to set deadline: %w", err)
		}
	}

	b.mu.RLock()
	tlsConn := utls.UClient(conn, b.getTLSConfig(serverName, nextProtos), b.clientHello)
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

func (b *BypassJA3Transport) dialUTLSContext(ctx context.Context, network, addr string, nextProtos []string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: 30 * time.Second}
	conn, err := dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, fmt.Errorf("tcp net dial fail: %w", err)
	}

	serverName := addr
	host, _, splitErr := net.SplitHostPort(addr)
	if splitErr == nil {
		serverName = host
	}

	tlsConn, err := b.tlsConnect(ctx, conn, serverName, nextProtos)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("tls connect fail: %w", err)
	}
	return tlsConn, nil
}

func (b *BypassJA3Transport) initTransports() {
	b.tr1 = *http.DefaultTransport.(*http.Transport).Clone()
	b.tr1.ForceAttemptHTTP2 = false
	b.tr1.MaxIdleConns = 1000
	b.tr1.MaxIdleConnsPerHost = 100
	b.tr1.IdleConnTimeout = 90 * time.Second
	b.tr1.TLSHandshakeTimeout = 10 * time.Second
	b.tr1.ExpectContinueTimeout = 1 * time.Second

	b.tr2 = http2.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return b.dialUTLSContext(ctx, network, addr, []string{"h2"})
		},
		ReadIdleTimeout: 30 * time.Second,
		PingTimeout:     15 * time.Second,
		IdleConnTimeout: 90 * time.Second,
	}
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
	helloID := utls.HelloChrome_Auto
	if disableHTTP2 {
		helloID = utls.HelloGolang
	}
	baseTransport := NewBypassJA3Transport(helloID, disableHTTP2)
	baseTransport.initTransports()
	return &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
		Transport: &BrowserHeaderRoundTripper{
			transport: baseTransport,
		},
	}
}
