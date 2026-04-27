package handlers

import (
	"context"

	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
)

// ScrapeTikTok implements [v1connect.RakerServerHandler].
func (r *RakerServer) ScrapeTikTok(context.Context, *v1.BinaryScrapeRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}