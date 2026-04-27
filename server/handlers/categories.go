package handlers

import (
	"context"

	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// EditCategory implements [v1connect.RakerServerHandler].
func (r *RakerServer) EditCategory(context.Context, *v1.EditCategoryRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateCategories implements [v1connect.RakerServerHandler].
func (r *RakerServer) UpdateCategories(context.Context, *v1.UpdateCategoriesRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}