package shared

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/imroc/req/v3"
)

var (
	Version        = "development"
	Hash           = "development"
	UserAgent      = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
	textareaRegExp = regexp.MustCompile(`<textarea class="form-control" rows="8">(.*?)</textarea>`)
	DefaultClient  = req.NewClient().ImpersonateChrome().SetUserAgent(UserAgent)
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
