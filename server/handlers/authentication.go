package handlers

import (
	"context"

	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SignInInstagram implements [v1connect.RakerServerHandler].
func (r *RakerServer) SignInInstagram(context.Context, *v1.SignUpRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// SignUpInstagram implements [v1connect.RakerServerHandler].
func (r *RakerServer) SignUpInstagram(context.Context, *v1.SignUpRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}