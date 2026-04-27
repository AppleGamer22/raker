package handlers

import (
	"context"

	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
)

// ScrapeVSCO implements [v1connect.RakerServerHandler].
func (r *RakerServer) ScrapeVSCO(context.Context, *v1.BinaryScrapeRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}