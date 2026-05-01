package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"connectrpc.com/connect"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/charmbracelet/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ScrapeVSCO implements [v1connect.RakerServerHandler].
func (server *RakerServer) ScrapeVSCO(ctx context.Context, request *v1.BinaryScrapeRequest) (*v1.ScrapeResponse, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	history, err := server.DBClient.HistoryGet(context.Background(), db.HistoryGetParams{
		PostType: db.PostTypeVsco,
		Post:     request.Post,
		Username: user.Username,
	})
	if err == nil {
		return &v1.ScrapeResponse{
			PostType:   v1.PostType_VSCO,
			PostOwner:  history.PostOwner,
			Post:       history.Post,
			PostDate:   timestamppb.New(history.PostDate),
			Files:      history.Files,
			Categories: history.Categories,
			Incognito:  history.Incognito,
		}, nil
	}

	URLs, username, cookies, err := shared.VSCO(request.Owner, request.Post)
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	localURLs := make([]string, 0, len(URLs))
	for _, urlString := range URLs {

		URL, err2 := url.Parse(urlString)
		if err2 != nil {
			err = errors.Join(err, err2)
			continue
		}
		if strings.Contains(urlString, ".ts") {
			localURLs = append(localURLs, fmt.Sprintf("%s.mp4", request.Post))
		} else if strings.Contains(URL.Path, "/poster/private") {
			localURLs = append(localURLs, fmt.Sprintf("%s.jpg", request.Post))
		} else {
			localURL := fmt.Sprintf("%s_%s", request.Post, path.Base(URL.Path))
			if !strings.HasSuffix(localURL, ".jpg") && !strings.HasSuffix(localURL, ".jpeg") {
				localURL = fmt.Sprintf("%s.jpeg", localURL)
			}
			localURLs = append(localURLs, localURL)
		}
	}
	localURLs, err2 := StorageHandler.SaveBundle(user, types.VSCO, username, localURLs, URLs, cookies)
	if err2 != nil {
		err = errors.Join(err, err2)
	}

	if len(localURLs) == 0 {
		return nil, connect.NewError(connect.CodeInternal, errors.New("no files could be saved"))
	}

	history, err2 = server.DBClient.HistoryAdd(context.Background(), db.HistoryAddParams{
		Username:   user.Username,
		PostType:   db.PostTypeSnapchat,
		PostOwner:  username,
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
