package handlers

import (
	"context"

	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
)

// ScrapeHighlight implements [v1connect.RakerServerHandler].
func (r *RakerServer) ScrapeHighlight(context.Context, *v1.UnaryScrapeRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}