package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var (
	threadsRegExpPostJSON = regexp.MustCompile(`<script type="application/json"\s*nonce="\S*"\s*data-content-len="\d*"\s*data-sjs>.*"data":{"media":(.*),"viewer":.*</script>`)
)

func Threads(owner, post string) ([]string, string, error) {
	htmlURL := fmt.Sprintf("https://www.threads.com/@%s/post/%s", owner, post)
	htmlRequest, err := http.NewRequest(http.MethodGet, htmlURL, nil)
	if err != nil {
		return []string{}, "", err
	}

	htmlRequest.Header.Add("Referer", "https://www.threads.com/")

	client := NewClient([]string{})

	htmlResponse, err := client.Do(htmlRequest)
	if err != nil {
		return []string{}, "", err
	}
	defer htmlResponse.Body.Close()

	htmlBody, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return []string{}, "", err
	}

	jsonMatch := threadsRegExpPostJSON.FindStringSubmatch(string(htmlBody))
	if len(jsonMatch) != 2 {
		// filename := fmt.Sprintf("instagram_%s.html", time.Now().Format("20060102_150405"))
		// if err := os.WriteFile(filename, htmlBody, 0644); err != nil {
		// 	fmt.Println("failed to save html body:", err)
		// } else {
		// 	fmt.Println("saved html body to", filename)
		// }
		// fmt.Println(string(htmlBody))
		return []string{}, "", fmt.Errorf("could not find JSON ID: %d matches", len(jsonMatch))
	}

	var threadsItem InstagramItem
	if err := json.Unmarshal([]byte(jsonMatch[1]), &threadsItem); err != nil {
		return []string{}, "", err
	}

	username := threadsItem.User.Username
	URLs := threadsItem.URLs()

	return URLs, username, nil
}
