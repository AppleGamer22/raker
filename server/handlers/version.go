package handlers

import (
	"fmt"
	"net/http"

	"github.com/AppleGamer22/rake/shared"
)

func Version(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, shared.Version)
}
