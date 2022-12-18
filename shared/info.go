package shared

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strings"
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
	response, err := http.Get("https://www.useragents.me/api")
	if err != nil {
		log.Println("could not retrieve the latest user agent")
		return
	}
	defer response.Body.Close()

	var data userAgentData
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil || len(data.Data) == 0 {
		log.Println("could not retrieve the latest user agent")
		return
	}

	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "mac"
	}

	for _, userAgent := range data.Data {
		if strings.Contains(strings.ToLower(userAgent.UserAgent), osName) {
			UserAgent = userAgent.UserAgent
			log.Println("using user agent", UserAgent)
			return
		}
	}
	log.Println("could not retrieve the latest user agent, using default")
}
