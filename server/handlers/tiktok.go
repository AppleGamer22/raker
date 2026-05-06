package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"

	"connectrpc.com/connect"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/charmbracelet/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ScrapeTikTok implements [v1connect.RakerServerHandler].
func (server *RakerServer) ScrapeTikTok(ctx context.Context, request *v1.BinaryScrapeRequest) (*v1.ScrapeResponse, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	history, err := server.DBClient.HistoryGet(context.Background(), db.HistoryGetParams{
		PostType: db.PostTypeTiktok,
		Post:     request.Post,
		Username: user.Username,
	})
	if err == nil {
		return &v1.ScrapeResponse{
			PostType:   v1.PostType_TikTok,
			PostOwner:  history.PostOwner,
			Post:       history.Post,
			PostDate:   timestamppb.New(history.PostDate),
			Files:      history.Files,
			Categories: history.Categories,
			Incognito:  history.Incognito,
		}, nil
	}

	tiktok := shared.NewTikTok(user.TiktokSessionID, user.TiktokSessionIDGuard)
	videoURLs, coverURLs, username, cookies, err := tiktok.Post(request.Owner, request.Post, request.GetIncognito())
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var videoURL string
	for i, urlString := range videoURLs {
		fileName := fmt.Sprintf("%s.mp4", request.Post)
		if err2 := StorageHandler.Save(user, db.PostTypeTiktok, username, fileName, urlString, cookies); err2 != nil {
			if i == len(videoURLs)-1 {
				err = errors.Join(err, err2)
			}
			continue
		}
		videoURL = fileName
		break
	}

	localURLs := make([]string, 0, len(coverURLs)+1)
	for _, urlString := range coverURLs {
		URL, err2 := url.Parse(urlString)
		if err2 != nil {
			err = errors.Join(err, err2)
			continue
		}
		fileName := fmt.Sprintf("%s_%s.jpeg", request.Post, path.Base(URL.Path))
		localURLs = append(localURLs, fileName)
	}

	localURLs, err2 := StorageHandler.SaveBundle(user, db.PostTypeTiktok, username, localURLs, coverURLs, cookies)
	if err2 != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Join(err, err2))
	}

	if len(localURLs) == 0 {
		return nil, connect.NewError(connect.CodeInternal, errors.New("no files could be saved"))
	}

	if len(videoURL) > 0 {
		localURLs = append([]string{videoURL}, localURLs...)
	}

	history, err2 = server.DBClient.HistoryAdd(context.Background(), db.HistoryAddParams{
		Username:   user.Username,
		PostType:   db.PostTypeTiktok,
		PostOwner:  username,
		Post:       request.Post,
		Files:      localURLs,
		Categories: []string{},
	})

	if err2 != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Join(err, err2))
	}

	return &v1.ScrapeResponse{
		PostType:   v1.PostType_TikTok,
		PostOwner:  history.PostOwner,
		Post:       history.Post,
		PostDate:   timestamppb.New(history.PostDate),
		Files:      history.Files,
		Categories: history.Categories,
		Incognito:  history.Incognito,
	}, nil
}
