package handlers

import (
	"context"

	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
)

// SearchHistory implements [v1connect.RakerServerHandler].
func (r *RakerServer) SearchHistory(context.Context, *v1.HistoryRequest) (*v1.HistoryResponse, error) {
	panic("unimplemented")
}

// SearchHistoryOwners implements [v1connect.RakerServerHandler].
func (r *RakerServer) SearchHistoryOwners(context.Context, *v1.HistoryRequest) (*v1.HistoryOwnersResponse, error) {
	panic("unimplemented")
}