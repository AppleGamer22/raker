package shared_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/AppleGamer22/raker/shared"
	"github.com/stretchr/testify/assert"
)

const postURL = "https://vsco.co/producedbymeee/media/68755110232dc1dc6ad53c48"

func getHTML(client *http.Client) (string, error) {
	htmlRequest, err := http.NewRequest(http.MethodGet, postURL, nil)
	if err != nil {
		return "", err
	}

	htmlRequest.Header.Set("User-Agent", shared.UserAgent)

	htmlResponse, err := client.Do(htmlRequest)
	if err != nil {
		return "", err
	}
	defer htmlResponse.Body.Close()

	body, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return "", err
	}

	if !strings.Contains(string(body), "<script>window.__PRELOADED_STATE__ =") {
		return string(body), errors.New("not found")
	}

	return string(body), nil
}

func testHTTP1(t *testing.T) {
	// client := http.DefaultClient
	client := shared.NewClient(true)
	errCount := 0
	for i := 0; i < 1e3; i++ {
		_, err := getHTML(client)
		// assert.Error(t, err, html)
		if err != nil {
			errCount++
		}
	}
	t.Log(errCount)
}

func testHTTP2(t *testing.T) {
	client := shared.NewClient(false)
	for i := 0; i < 1e3; i++ {
		_, err := getHTML(client)
		assert.NoError(t, err)
	}
}

func TestClientHello(t *testing.T) {
	t.Run("HTTP/1.1", testHTTP1)
	t.Run("HTTP/2", testHTTP2)
}
