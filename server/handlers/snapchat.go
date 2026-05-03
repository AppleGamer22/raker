package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"connectrpc.com/connect"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/charmbracelet/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ScrapeInstagram implements [v1connect.RakerServerHandler].
func (server *RakerServer) ScrapeSnapchat(ctx context.Context, request *v1.BinaryScrapeRequest) (*v1.ScrapeResponse, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	history, err := server.DBClient.HistoryGet(context.Background(), db.HistoryGetParams{
		PostType: db.PostTypeSnapchat,
		Post:     request.Post,
		Username: user.Username,
	})
	if err == nil {
		return &v1.ScrapeResponse{
			PostType:   v1.PostType_Snapchat,
			PostOwner:  history.PostOwner,
			Post:       history.Post,
			PostDate:   timestamppb.New(history.PostDate),
			Files:      history.Files,
			Categories: history.Categories,
			Incognito:  history.Incognito,
		}, nil
	}

	result, _, err := shared.Snapchat(request.Owner, request.Post)
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	localURLs := make([]string, 0, len(result.URLs))
	remoteURLs := make([]string, 0, len(result.URLs))
	for _, u := range result.URLs {
		URL, err2 := url.Parse(u.URL)
		if err2 != nil {
			err = errors.Join(err, err2)
			continue
		}
		fileName := fmt.Sprintf("%s_%s", request.Post, path.Base(URL.Path))
		if u.IsVideo {
			fileName = fmt.Sprintf("%s.mp4", fileName)
		} else {
			fileName = fmt.Sprintf("%s.jpg", fileName)
		}
		localURLs = append(localURLs, fileName)
		remoteURLs = append(remoteURLs, u.URL)
	}

	localURLs, err2 := StorageHandler.SaveBundle(user, db.PostTypeSnapchat, result.Username, localURLs, remoteURLs, []*http.Cookie{})
	if err2 != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Join(err, err2))
	}

	if len(localURLs) == 0 {
		return nil, connect.NewError(connect.CodeInternal, errors.New("no files could be saved"))
	}

	history, err2 = server.DBClient.HistoryAdd(context.Background(), db.HistoryAddParams{
		Username:   user.Username,
		PostType:   db.PostTypeSnapchat,
		PostOwner:  result.Username,
		Post:       request.Post,
		Files:      localURLs,
		Categories: []string{},
	})

	if err2 != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Join(err, err2))
	}

	return &v1.ScrapeResponse{
		PostType:   v1.PostType_Snapchat,
		PostOwner:  history.PostOwner,
		Post:       history.Post,
		PostDate:   timestamppb.New(history.PostDate),
		Files:      history.Files,
		Categories: history.Categories,
		Incognito:  history.Incognito,
	}, nil
}
