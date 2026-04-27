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

var unauthenticatedProcedures = map[string]struct{}{
	"/raker.v1.RakerServer/SignInInstagram": {},
	"/raker.v1.RakerServer/SignUpInstagram": {},
}

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
		HttpOnly: false,
	}

	// Based on https://connectrpc.com/docs/go/headers-and-trailers/#headers
	callInfo, ok := connect.CallInfoForHandlerContext(ctx)
	if !ok {
		return nil, errors.New("can't access headers: no CallInfo for handler context")
	}

	callInfo.ResponseHeader().Set("Set-Cookie", cookie.String())

	return &emptypb.Empty{}, nil
}

func (server *RakerServer) GetUserFromCookie(cookie *http.Cookie) (db.User, error) {
	username, err := server.Authenticator.Parse(cookie.Value)
	if err != nil {
		return db.User{}, err
	}

	user, err := server.DBClient.UserGet(context.Background(), username)
	return user, err
}

func (server *RakerServer) NewAuthInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if _, ok := unauthenticatedProcedures[req.Spec().Procedure]; ok {
				// sign-in/up
				return next(ctx, req)
			}

			cookies, err := http.ParseCookie(req.Header().Get("Cookie"))

			if err != nil {
				log.Error(err)
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("no token provided"),
				)
			}

			for _, cookie := range cookies {
				if cookie.Name != "jwt" {
					continue
				}

				user, err := server.GetUserFromCookie(cookie)
				if err != nil {
					log.Error(err)
					return nil, err
				}

				ctxWithUser := context.WithValue(ctx, authenticatedUserKey, user)
				return next(ctxWithUser, req)
			}

			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("no token provided"),
			)
		}
	}
}

// EditUserCredentials implements [v1connect.RakerServerHandler].
func (server *RakerServer) EditUserCredentials(ctx context.Context, request *v1.EditUserCredentialsRequest) (*emptypb.Empty, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}
	passwordHash := user.PasswordHash
	shouldUpdateHash := request.Password != nil && len(*request.Password) > 0
	if shouldUpdateHash {
		hashed, err := authenticator.Hash(*request.Password)
		if err != nil {
			log.Error(err)
			return &emptypb.Empty{}, connect.NewError(connect.CodeInternal, errors.New("failed to store credentials securely"))
		}
		passwordHash = hashed
	}

	sessionID := user.InstagramSessionID
	shouldUpdateSessionID := request.SessionId != nil && len(*request.SessionId) > 0
	if shouldUpdateSessionID {
		sessionID = *request.SessionId
	}

	userID := user.InstagramUserID
	shouldUpdateUserID := request.UserId != nil && len(*request.UserId) > 0
	if shouldUpdateUserID {
		sessionID = *request.UserId
	}

	if shouldUpdateHash || shouldUpdateSessionID || shouldUpdateUserID {
		err := server.DBClient.UserUpdateInstagramSession(ctx, db.UserUpdateInstagramSessionParams{
			Username:           user.Username,
			PasswordHash:       passwordHash,
			InstagramSessionID: sessionID,
			InstagramUserID:    userID,
		})
		if err != nil {
			log.Error(err)
			return &emptypb.Empty{}, connect.NewError(connect.CodeInternal, err)
		}
	}

	return &emptypb.Empty{}, nil
}
