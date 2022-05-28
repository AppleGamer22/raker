package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AppleGamer22/rake/shared"
)

type version struct {
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

func Information(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Query().Get("about") {
	case "":
		jsonPayload := version{shared.Version, shared.Hash}
		data, err := json.Marshal(jsonPayload)
		if err != nil {
			http.Error(writer, "could not process version information", http.StatusInternalServerError)
		}
		fmt.Fprint(writer, string(data))
	case "version":
		fmt.Fprint(writer, shared.Version)
	case "hash":
		fmt.Fprint(writer, shared.Hash)
	default:
		http.Error(writer, `about query parameter can be either "", "version" or "hash"`, http.StatusBadRequest)
	}
}
