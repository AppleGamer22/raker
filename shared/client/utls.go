package client

// TODO: look into this Claude Sonnet 4 output

// import (
// 	"crypto/tls"
// 	"fmt"
// 	"io"
// 	"net"
// 	"net/http"
// 	"time"

// 	utls "github.com/refraction-networking/utls"
// 	"golang.org/x/net/http2"
// )

// // UTLSRoundTripper wraps http2.Transport with utls for custom TLS fingerprinting
// type UTLSRoundTripper struct {
// 	clientHelloID utls.ClientHelloID
// 	transport     *http2.Transport
// }

// // NewUTLSRoundTripper creates a new round tripper with utls support
// func NewUTLSRoundTripper(clientHelloID utls.ClientHelloID) *UTLSRoundTripper {
// 	return &UTLSRoundTripper{
// 		clientHelloID: clientHelloID,
// 		transport: &http2.Transport{
// 			TLSClientConfig: &tls.Config{
// 				InsecureSkipVerify: false,
// 			},
// 			// Allow HTTP/2 connection reuse
// 			AllowHTTP: false,
// 			// Set reasonable timeouts
// 			ReadIdleTimeout: 30 * time.Second,
// 			PingTimeout:     15 * time.Second,
// 		},
// 	}
// }

// // dialTLSWithUTLS performs TLS handshake using utls
// func (rt *UTLSRoundTripper) dialTLSWithUTLS(network, addr string) (net.Conn, error) {
// 	// Establish TCP connection
// 	conn, err := net.DialTimeout(network, addr, 30*time.Second)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get hostname for SNI
// 	host, _, err := net.SplitHostPort(addr)
// 	if err != nil {
// 		conn.Close()
// 		return nil, err
// 	}

// 	// Create utls connection with specified fingerprint
// 	utlsConn := utls.UClient(conn, &utls.Config{
// 		ServerName:         host,
// 		InsecureSkipVerify: false,
// 	}, rt.clientHelloID)

// 	// Perform TLS handshake
// 	err = utlsConn.Handshake()
// 	if err != nil {
// 		conn.Close()
// 		return nil, err
// 	}

// 	return utlsConn, nil
// }

// // RoundTrip implements http.RoundTripper interface
// func (rt *UTLSRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
// 	// Set up custom dialer for the transport
// 	rt.transport.DialTLS = rt.dialTLSWithUTLS

// 	return rt.transport.RoundTrip(req)
// }

// // CreateHTTP2ClientWithUTLS creates an HTTP client configured for HTTP/2 with utls
// func CreateHTTP2ClientWithUTLS(clientHelloID utls.ClientHelloID) *http.Client {
// 	return &http.Client{
// 		Transport: NewUTLSRoundTripper(clientHelloID),
// 		Timeout:   60 * time.Second,
// 	}
// }

// // Example usage function
// func main() {
// 	// Create client with Chrome fingerprint
// 	client := CreateHTTP2ClientWithUTLS(utls.HelloChrome_Auto)

// 	// Make HTTP/2 request
// 	resp, err := client.Get("https://httpbin.org/get")
// 	if err != nil {
// 		fmt.Printf("Error making request: %v\n", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Check protocol version
// 	fmt.Printf("Protocol: %s\n", resp.Proto)
// 	fmt.Printf("Status: %s\n", resp.Status)

// 	// Read response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Printf("Error reading response: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("Response body length: %d bytes\n", len(body))

// 	// Example with different fingerprints
// 	testDifferentFingerprints()
// }

// // testDifferentFingerprints demonstrates using various TLS fingerprints
// func testDifferentFingerprints() {
// 	fingerprints := []struct {
// 		name string
// 		id   utls.ClientHelloID
// 	}{
// 		{"Chrome", utls.HelloChrome_Auto},
// 		{"Firefox", utls.HelloFirefox_Auto},
// 		{"Safari", utls.HelloSafari_Auto},
// 		{"Edge", utls.HelloEdge_Auto},
// 	}

// 	for _, fp := range fingerprints {
// 		fmt.Printf("\nTesting with %s fingerprint:\n", fp.name)

// 		client := CreateHTTP2ClientWithUTLS(fp.id)

// 		resp, err := client.Get("https://httpbin.org/headers")
// 		if err != nil {
// 			fmt.Printf("Error with %s: %v\n", fp.name, err)
// 			continue
// 		}

// 		fmt.Printf("Protocol: %s, Status: %s\n", resp.Proto, resp.Status)
// 		resp.Body.Close()
// 	}
// }

// // AdvancedUTLSClient provides more configuration options
// type AdvancedUTLSClient struct {
// 	client        *http.Client
// 	clientHelloID utls.ClientHelloID
// 	timeout       time.Duration
// }

// // NewAdvancedUTLSClient creates a configurable utls HTTP client
// func NewAdvancedUTLSClient(clientHelloID utls.ClientHelloID, timeout time.Duration) *AdvancedUTLSClient {
// 	transport := &http2.Transport{
// 		TLSClientConfig: &tls.Config{
// 			InsecureSkipVerify: false,
// 			MinVersion:         tls.VersionTLS12,
// 		},
// 		AllowHTTP:       false,
// 		ReadIdleTimeout: 30 * time.Second,
// 		PingTimeout:     15 * time.Second,
// 	}

// 	// Custom dialer with utls
// 	transport.DialTLS = func(network, addr string) (net.Conn, error) {
// 		conn, err := net.DialTimeout(network, addr, 30*time.Second)
// 		if err != nil {
// 			return nil, err
// 		}

// 		host, _, err := net.SplitHostPort(addr)
// 		if err != nil {
// 			conn.Close()
// 			return nil, err
// 		}

// 		utlsConn := utls.UClient(conn, &utls.Config{
// 			ServerName:         host,
// 			InsecureSkipVerify: false,
// 		}, clientHelloID)

// 		err = utlsConn.Handshake()
// 		if err != nil {
// 			conn.Close()
// 			return nil, err
// 		}

// 		return utlsConn, nil
// 	}

// 	return &AdvancedUTLSClient{
// 		client: &http.Client{
// 			Transport: transport,
// 			Timeout:   timeout,
// 		},
// 		clientHelloID: clientHelloID,
// 		timeout:       timeout,
// 	}
// }
