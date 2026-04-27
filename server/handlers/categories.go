package handlers

import (
	"context"
	"errors"
	"slices"

	"connectrpc.com/connect"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/charmbracelet/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *RakerServer) updateCategory(ctx context.Context, user db.User, oldCategoryName, newCategoryName string) error {
	tx, err := server.DBConnection.Begin()
	if err != nil {
		log.Error(err, "new_category", newCategoryName, "old_category", oldCategoryName)
		return err
	}
	defer tx.Rollback()

	qtx := server.DBClient.WithTx(tx)
	// Rename: update all user's histories with the old category, then update user's list
	err = qtx.HistoriesCategoryRename(ctx, db.HistoriesCategoryRenameParams{
		Username:    user.Username,
		OldCategory: oldCategoryName,
		NewCategory: newCategoryName,
	})
	if err != nil {
		log.Error(err, "new_category", newCategoryName, "old_category", oldCategoryName)
		return err
	}

	err = qtx.UserCategoryRemove(ctx, db.UserCategoryRemoveParams{
		Username: user.Username,
		Category: oldCategoryName,
	})
	if err != nil {
		log.Error(err, "new_category", newCategoryName, "old_category", oldCategoryName)
		return err
	}

	err = qtx.UserCategoryAdd(ctx, db.UserCategoryAddParams{
		Username: user.Username,
		Category: newCategoryName,
	})
	if err != nil {
		log.Error(err, "new_category", newCategoryName, "old_category", oldCategoryName)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error(err, "new_category", newCategoryName, "old_category", oldCategoryName)
		return err
	}

	return nil
}

// EditCategory implements [v1connect.RakerServerHandler].
func (server *RakerServer) EditCategory(ctx context.Context, request *v1.EditCategoryRequest) (*emptypb.Empty, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	oldCategoryName := request.OldCategory
	newCategoryName := request.NewCategory

	if len(newCategoryName) == 0 {
		// delete category
		err := server.DBClient.UserCategoryRemove(ctx, db.UserCategoryRemoveParams{
			Username: user.Username,
			Category: oldCategoryName,
		})
		if err != nil {
			log.Error(err, "category", oldCategoryName)
			return &emptypb.Empty{}, connect.NewError(connect.CodeInternal, err)
		}
		return &emptypb.Empty{}, nil
	}

	if slices.Contains(user.Categories, oldCategoryName) && !slices.Contains(user.Categories, newCategoryName) {
		// update category
		if err := server.updateCategory(ctx, user, oldCategoryName, newCategoryName); err != nil {
			log.Error(err, "new_category", newCategoryName, "old_category", oldCategoryName)
			return &emptypb.Empty{}, connect.NewError(connect.CodeInternal, err)
		}
	} else if oldCategoryName == newCategoryName && !slices.Contains(user.Categories, newCategoryName) {
		// add category
		err := server.DBClient.UserCategoryAdd(ctx, db.UserCategoryAddParams{
			Username: user.Username,
			Category: newCategoryName,
		})
		if err != nil {
			log.Error(err, "new_category", newCategoryName, "old_category", oldCategoryName)
			return &emptypb.Empty{}, connect.NewError(connect.CodeInternal, err)
		}
	}

	return &emptypb.Empty{}, nil
}

// UpdateCategories implements [v1connect.RakerServerHandler].
func (server *RakerServer) UpdateCategories(ctx context.Context, request *v1.UpdateCategoriesRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}
