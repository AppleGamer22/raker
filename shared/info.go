package shared

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	utls "github.com/refraction-networking/utls"
)

var (
	Version        = "development"
	Hash           = "development"
	UserAgent      = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
	textareaRegExp = regexp.MustCompile(`<textarea class="form-control" rows="8">(.*?)</textarea>`)
)

type userAgentData struct {
	UserAgent string `json:"ua"`
}

func init() {
	log.SetReportCaller(true)
	log.SetTimeFormat(time.RFC3339)
	log.SetLevel(log.DebugLevel)
	logger := slog.New(log.Default())
	slog.SetDefault(logger)
}

// CustomDialTLS uses uTLS for TLS handshake
func CustomDialTLS(ctx context.Context, network, addr string) (conn *utls.UConn, err error) {
	colonPos := strings.LastIndex(addr, ":")
	if colonPos == -1 {
		colonPos = len(addr)
	}
	hostname := addr[:colonPos]
	config := &utls.Config{
		ServerName: hostname,
		NextProtos: []string{"h2", "http/1.1"},
	}
	// Standard TCP connection
	tcpConn, err := (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// Pick a ClientHelloID (imitate Chrome, Firefox, etc.)
	uconn := utls.UClient(tcpConn, config, utls.HelloChrome_Auto)
	err = uconn.Handshake()
	if err != nil {
		return nil, err
	}
	return uconn, nil
}

var DefaultClient = &http.Client{
	Transport: &http.Transport{
		// ForceAttemptHTTP2: true,
		// TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return CustomDialTLS(ctx, network, addr)
		},
	},
}

func init() {
	response, err := http.Get("https://www.useragents.me/")
	if err != nil {
		log.Error("could not retrieve the latest user agent")
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return
	}

	textarea := textareaRegExp.FindAllStringSubmatch(string(body), 1)
	if len(textarea) == 0 {
		log.Error("couldn't find latest user agent string")
		return
	}

	var data []userAgentData
	if err := json.Unmarshal([]byte(textarea[0][1]), &data); err != nil || len(data) == 0 {
		log.Error("could not retrieve the latest user agent", "err", err)
		return
	}

	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "mac"
	}

	for _, userAgent := range data {
		if strings.Contains(strings.ToLower(userAgent.UserAgent), osName) && !strings.Contains(userAgent.UserAgent, "Samsung") {
			UserAgent = userAgent.UserAgent
			log.Infof("using user agent %s", UserAgent)
			return
		}
	}
	log.Warn("could not retrieve the latest user agent, using default")
}
