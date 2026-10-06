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

func getHTML(client *http.Client) (string, int, error) {
	htmlRequest, err := http.NewRequest(http.MethodGet, postURL, nil)
	if err != nil {
		return "", 0, err
	}

	htmlRequest.Header.Set("User-Agent", shared.UserAgent)

	htmlResponse, err := client.Do(htmlRequest)
	if err != nil {
		return "", 0, err
	}
	defer htmlResponse.Body.Close()

	body, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return "", htmlResponse.ProtoMajor, err
	}

	if !strings.Contains(string(body), "<script>window.__PRELOADED_STATE__ =") {
		return string(body), htmlResponse.ProtoMajor, errors.New("not found")
	}

	return string(body), htmlResponse.ProtoMajor, nil
}

func testHTTP1(t *testing.T) {
	// client := http.DefaultClient
	client := shared.NewClient([]string{"http/1.1"})
	errCount := 0
	for i := 0; i < 1e0; i++ {
		_, protoMajor, err := getHTML(client)
		assert.Error(t, err)
		assert.Equal(t, 1, protoMajor)
		if err != nil {
			errCount++
		}
	}
	t.Log(errCount)
}

func testHTTP2(t *testing.T) {
	client := shared.NewClient([]string{"h2"})
	for i := 0; i < 1e0; i++ {
		_, protoMajor, err := getHTML(client)
		assert.NoError(t, err)
		assert.Equal(t, 2, protoMajor)
	}
}

func TestClientHello(t *testing.T) {
	t.Run("HTTP/1.1", testHTTP1)
	t.Run("HTTP/2", testHTTP2)
}
