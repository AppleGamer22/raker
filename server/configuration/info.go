package configuration

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AppleGamer22/raker/shared"
)

type version struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

func Information(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Query().Get("about") {
	case "":
		jsonPayload := version{"raker", shared.Version, shared.Hash}
		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(jsonPayload); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	case "version":
		fmt.Fprint(writer, shared.Version)
	case "hash":
		fmt.Fprint(writer, shared.Hash)
	default:
		http.Error(writer, `about query parameter can be either "", "version" or "hash"`, http.StatusBadRequest)
	}
}
