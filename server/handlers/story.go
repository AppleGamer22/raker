package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"

	"connectrpc.com/connect"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ScrapeStory implements [v1connect.RakerServerHandler].
func (server *RakerServer) ScrapeStory(ctx context.Context, request *v1.UnaryScrapeRequest) (*v1.ScrapeResponse, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	history, err := server.DBClient.HistoryGet(context.Background(), db.HistoryGetParams{
		PostType: db.PostTypeStory,
		Post:     request.Post,
		Username: user.Username,
	})
	if err == nil {
		return &v1.ScrapeResponse{
			PostType:   v1.PostType_Story,
			PostOwner:  history.PostOwner,
			Post:       history.Post,
			PostDate:   timestamppb.New(history.PostDate),
			Files:      history.Files,
			Categories: history.Categories,
			Incognito:  history.Incognito,
		}, nil
	}

	instagram := shared.NewInstagram(user.InstagramSessionID, user.InstagramUserID)
	URLs, username, err := instagram.Story(request.Post)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	newURLs := make([]string, 0, len(URLs))
	localURLs := make([]string, 0, len(URLs))
	for _, urlString := range URLs {
		URL, err2 := url.Parse(urlString)
		if err2 != nil {
			err = errors.Join(err, err2)
			continue
		}

		fileName := path.Base(URL.Path)
		count, err2 := server.DBClient.HistoryCountByFile(context.Background(), db.HistoryCountByFileParams{
			PostType:  db.PostTypeStory,
			PostOwner: username,
			File:      fileName,
			Username:  user.Username,
		})
		if err2 != nil || count > 0 {
			continue
		}

		localURLs = append(localURLs, fileName)
		newURLs = append(newURLs, urlString)
	}

	localURLs, err2 := StorageHandler.SaveBundle(user, db.PostTypeStory, username, localURLs, newURLs, []*http.Cookie{})
	if err2 != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Join(err, err2))
	}

	if len(localURLs) == 0 {
		return nil, connect.NewError(connect.CodeInternal, errors.New("no files could be saved"))
	}

	history, err2 = server.DBClient.HistoryAdd(context.Background(), db.HistoryAddParams{
		Username:   user.Username,
		PostType:   db.PostTypeStory,
		PostOwner:  username,
		Post:       uuid.NewString(),
		Files:      localURLs,
		Categories: []string{},
	})

	if err2 != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Join(err, err2))
	}

	return &v1.ScrapeResponse{
		PostType:   v1.PostType_Story,
		PostOwner:  history.PostOwner,
		Post:       history.Post,
		PostDate:   timestamppb.New(history.PostDate),
		Files:      history.Files,
		Categories: history.Categories,
		Incognito:  history.Incognito,
	}, nil
}
