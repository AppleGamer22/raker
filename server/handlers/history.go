package handlers

import (
	"context"
	"errors"
	"math"
	"sort"

	"connectrpc.com/connect"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/db"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PostTypeDB2PB(pt db.PostType) v1.PostType {
	switch pt {
	case db.PostTypeInstagram:
		return v1.PostType_Instagram
	case db.PostTypeHighlight:
		return v1.PostType_Highlight
	case db.PostTypeStory:
		return v1.PostType_Story
	case db.PostTypeTiktok:
		return v1.PostType_TikTok
	case db.PostTypeSnapchat:
		return v1.PostType_Snapchat
	case db.PostTypeVsco:
		return v1.PostType_VSCO
	default:
		return v1.PostType_Instagram
	}
}

func PostTypePB2DB(pt v1.PostType) db.PostType {
	switch pt {
	case v1.PostType_Instagram:
		return db.PostTypeInstagram
	case v1.PostType_Highlight:
		return db.PostTypeHighlight
	case v1.PostType_Story:
		return db.PostTypeStory
	case v1.PostType_TikTok:
		return db.PostTypeTiktok
	case v1.PostType_Snapchat:
		return db.PostTypeSnapchat
	case v1.PostType_VSCO:
		return db.PostTypeVsco
	default:
		return db.PostTypeInstagram
	}
}

func HistoryToScrapeResponse(user db.User, history db.History) *v1.ScrapeResponse {
	return resolveVSCOMetadata(user, &v1.ScrapeResponse{
		PostOwner:  history.PostOwner,
		Post:       history.Post,
		Files:      history.Files,
		PostDate:   timestamppb.New(history.PostDate),
		Categories: history.Categories,
		Incognito:  history.Incognito,
		PostType:   PostTypeDB2PB(history.PostType),
	})
}

// SearchHistory implements [v1connect.RakerServerHandler].
func (server *RakerServer) SearchHistory(ctx context.Context, request *v1.HistoryRequest) (*v1.HistoryResponse, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	postTypes := make([]db.PostType, 0, len(request.Types))
	for _, postType := range request.Types {
		postTypes = append(postTypes, PostTypePB2DB(postType))
	}

	count, err := server.DBClient.HistoryCount(context.Background(), db.HistoryCountParams{
		PostTypes:  postTypes,
		Exclusive:  request.Exclusive,
		Categories: request.Categories,
		PostOwners: request.Owners,
		Username:   user.Username,
	})
	if err != nil {
		return &v1.HistoryResponse{}, connect.NewError(connect.CodeInternal, err)
	}

	page := request.Page
	pages := int64(math.Ceil(float64(count) / float64(request.PageSize)))
	if pages == 0 {
		page = 1
		pages = 1
	} else if page > pages {
		page = pages
	}

	if request.Exclusive {
		sort.Strings(request.Categories)
	}

	histories, err := server.DBClient.HistoryGetPage(context.Background(), db.HistoryGetPageParams{
		PostTypes:  postTypes,
		Exclusive:  request.Exclusive,
		Categories: request.Categories,
		PostOwners: request.Owners,
		Username:   user.Username,
		Page:       int32((page - 1) * int64(request.PageSize)),
		PageSize:   30,
	})

	if err != nil {
		return &v1.HistoryResponse{}, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.HistoryResponse{
		TotalCount: count,
		Histories: (func() []*v1.ScrapeResponse {
			output := make([]*v1.ScrapeResponse, 0, len(histories))
			for _, history := range histories {
				output = append(output, HistoryToScrapeResponse(user, history))
			}
			return output
		})(),
	}, nil
}

// SearchHistoryOwners implements [v1connect.RakerServerHandler].
func (server *RakerServer) SearchHistoryOwners(ctx context.Context, request *v1.HistoryOwnersRequest) (*v1.HistoryOwnersResponse, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	postTypes := make([]db.PostType, 0, len(request.Types))
	for _, postType := range request.Types {
		postTypes = append(postTypes, PostTypePB2DB(postType))
	}

	result, err := server.DBClient.HistoryOwners(ctx, db.HistoryOwnersParams{
		PostTypes:  postTypes,
		Exclusive:  request.Exclusive,
		Categories: request.Categories,
		PostOwner:  request.Owner,
		Username:   user.Username,
	})

	if err != nil {
		return &v1.HistoryOwnersResponse{}, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.HistoryOwnersResponse{
		Owners: func() []*v1.HistoryOwnersResponse_HistoryOwner {
			owners := make([]*v1.HistoryOwnersResponse_HistoryOwner, 0, len(result))
			for _, historyOwner := range result {
				owners = append(owners, &v1.HistoryOwnersResponse_HistoryOwner{
					Owner: historyOwner.PostOwner,
					Type:  PostTypeDB2PB(historyOwner.PostType),
				})
			}
			return owners
		}(),
	}, nil
}
