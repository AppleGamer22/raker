package shared

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

const (
	Version = "development"
	Hash    = "development"
)

var UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

type userAgentData struct {
	Data []struct {
		UserAgent string `json:"ua"`
	} `json:"data"`
}

func init() {
	log.SetReportCaller(true)
	log.SetTimeFormat(time.RFC3339)
	log.SetLevel(log.DebugLevel)
}

func init() {
	response, err := http.Get("https://www.useragents.me/api")
	if err != nil {
		log.Error("could not retrieve the latest user agent")
		return
	}
	defer response.Body.Close()

	var data userAgentData
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil || len(data.Data) == 0 {
		log.Error("could not retrieve the latest user agent", "err", err)
		return
	}

	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "mac"
	}

	for _, userAgent := range data.Data {
		if strings.Contains(strings.ToLower(userAgent.UserAgent), osName) {
			UserAgent = userAgent.UserAgent
			log.Infof("using user agent %s", UserAgent)
			return
		}
	}
	log.Warn("could not retrieve the latest user agent, using default")
}
