package handlers

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"github.com/AppleGamer22/raker/server/authenticator"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/charmbracelet/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SignUpInstagram implements [v1connect.RakerServerHandler].
func (server *RakerServer) SignUpInstagram(ctx context.Context, request *v1.SignUpRequest) (*emptypb.Empty, error) {
	username := request.Username
	password := request.Password
	sessionID := ""
	userID := ""
	if request.SessionId != nil && request.UserId != nil {
		sessionID = *request.SessionId
		userID = *request.UserId
	}
	_, err := server.DBClient.UserGet(context.Background(), username)
	if err == nil {
		log.Error("username already exists", "username", username)
		return &emptypb.Empty{}, connect.NewError(connect.CodeAlreadyExists, errors.New("username already exists"))
	}

	hashed, err := authenticator.Hash(password)
	if err != nil {
		return &emptypb.Empty{}, connect.NewError(connect.CodeCanceled, errors.New("failed to store credentials securely"))
	}

	err = server.DBClient.UserAdd(context.Background(), db.UserAddParams{
		Username:           username,
		PasswordHash:       hashed,
		InstagramSessionID: sessionID,
		InstagramUserID:    userID,
	})
	if err != nil {
		log.Error(err)
		return &emptypb.Empty{}, connect.NewError(connect.CodeInternal, err)
	}

	_, err = server.SignInInstagram(ctx, &v1.SignInRequest{
		Username: request.Username,
		Password: request.Password,
	})

	return &emptypb.Empty{}, err
}

// SignInInstagram implements [v1connect.RakerServerHandler].
func (server *RakerServer) SignInInstagram(ctx context.Context, request *v1.SignInRequest) (*emptypb.Empty, error) {
	username := request.Username
	password := request.Password
	user, err := server.DBClient.UserGet(context.Background(), username)
	if err != nil {
		log.Error(err)
		return &emptypb.Empty{}, connect.NewError(connect.CodeUnauthenticated, errors.New("incorrect credentials"))
	}

	if err := authenticator.Compare(user.PasswordHash, password); err != nil {
		log.Error(err)
		return &emptypb.Empty{}, connect.NewError(connect.CodeUnauthenticated, errors.New("incorrect credentials"))

	}

	webToken, expiry, err := server.Authenticator.Sign(user.Username)
	if err != nil {
		log.Error(err, "ID", user.Username)
		return &emptypb.Empty{}, connect.NewError(connect.CodeUnauthenticated, errors.New("incorrect credentials"))
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    webToken,
		Path:     "/",
		Expires:  expiry,
		Secure:   server.Configuration.SecureCookie,
		HttpOnly: true,
	}

	// Based on https://connectrpc.com/docs/go/headers-and-trailers/#headers
	callInfo, ok := connect.CallInfoForHandlerContext(ctx)
	if !ok {
		return nil, errors.New("can't access headers: no CallInfo for handler context")
	}

	callInfo.ResponseHeader().Set("Set-Cookie", cookie.String())

	return &emptypb.Empty{}, nil
}
