package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/AppleGamer22/raker/server/authenticator"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/buf/proto/raker/v1/passkey"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/charmbracelet/log"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

var unauthenticatedProcedures = map[string]struct{}{
	"/raker.v1.RakerServer/SignInInstagram": {},
	"/raker.v1.RakerServer/SignUpInstagram": {},
	"/raker.v1.RakerServer/BeginSignUp":     {},
	"/raker.v1.RakerServer/FinishSignUp":    {},
	"/raker.v1.RakerServer/BeginSignIn":     {},
	"/raker.v1.RakerServer/FinishSignIn":    {},
}

func (server *RakerServer) ApproveSession(ctx context.Context, user db.User) error {
	webToken, expiry, err := server.Authenticator.Sign(user.Username)
	if err != nil {
		log.Error(err, "ID", user.Username)
		return connect.NewError(connect.CodeUnauthenticated, errors.New("incorrect credentials"))
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
		return errors.New("can't access headers: no CallInfo for handler context")
	}

	callInfo.ResponseHeader().Set("Set-Cookie", cookie.String())
	return nil
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
	_, err := server.DBClient.UserGetByUsername(context.Background(), username)
	if err == nil {
		log.Error("username already exists", "username", username)
		return &emptypb.Empty{}, connect.NewError(connect.CodeAlreadyExists, errors.New("username already exists"))
	}

	hashed, err := authenticator.Hash(password)
	if err != nil {
		return &emptypb.Empty{}, connect.NewError(connect.CodeCanceled, errors.New("failed to store credentials securely"))
	}

	_, err = server.DBClient.UserAdd(context.Background(), db.UserAddParams{
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
	user, err := server.DBClient.UserGetByUsername(context.Background(), username)
	if err != nil {
		log.Error(err)
		return &emptypb.Empty{}, connect.NewError(connect.CodeUnauthenticated, errors.New("incorrect credentials"))
	}

	if err := authenticator.Compare(user.PasswordHash, password); err != nil {
		log.Error(err)
		return &emptypb.Empty{}, connect.NewError(connect.CodeUnauthenticated, errors.New("incorrect credentials"))

	}

	if err := server.ApproveSession(ctx, user); err != nil {
		log.Error(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (server *RakerServer) GetUserFromCookie(cookie *http.Cookie) (db.User, error) {
	username, err := server.Authenticator.Parse(cookie.Value)
	if err != nil {
		return db.User{}, err
	}

	user, err := server.DBClient.UserGetByUsername(context.Background(), username)
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

	fmt.Printf("%+v\n", request)

	sessionID := user.InstagramSessionID
	shouldUpdateSessionID := request.SessionId != nil && len(*request.SessionId) > 0
	if shouldUpdateSessionID {
		sessionID = *request.SessionId
	}

	userID := user.InstagramUserID
	shouldUpdateUserID := request.UserId != nil && len(*request.UserId) > 0
	if shouldUpdateUserID {
		userID = *request.UserId
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

// BeginSignUp implements [v1connect.RakerServerHandler].
func (server *RakerServer) BeginSignUp(ctx context.Context, request *passkey.BeginSignUpRequest) (*passkey.BeginSignUpResponse, error) {
	if request.Username == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username required"))
	}

	_, err := server.DBClient.UserGetByUsername(context.Background(), request.Username)
	if err == nil {
		log.Error("username already exists", "username", request.Username)
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("username already exists"))
	}

	userEntity := &authenticator.UserEntity{ID: uuid.New(), Username: request.Username}
	options, sessionData, err := server.WebAuthn.BeginRegistration(userEntity)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	sessionID := uuid.NewString()
	server.WebAuthnSessionStore.Set(sessionID, authenticator.WebAuthnSession{
		WebAuthnSession: *sessionData,
		UserID:          userEntity.ID,
		Username:        userEntity.Username,
	})

	optionsJSON, _ := json.Marshal(options)

	return &passkey.BeginSignUpResponse{
		SessionId:   sessionID,
		OptionsJson: string(optionsJSON),
	}, nil
}

// FinishSignUp implements [v1connect.RakerServerHandler].
func (server *RakerServer) FinishSignUp(ctx context.Context, request *passkey.FinishSignUpRequest) (*emptypb.Empty, error) {
	sessionData, ok := server.WebAuthnSessionStore.Get(request.GetSessionId())
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("session expired or invalid"))
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "", bytes.NewReader([]byte(request.GetResponseJson())))
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	userEntity := &authenticator.UserEntity{ID: sessionData.UserID, Username: sessionData.Username}
	credential, err := server.WebAuthn.FinishRegistration(userEntity, sessionData.WebAuthnSession, httpReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("passkey verification failed: %w", err))
	}

	var transports []string
	for _, t := range credential.Transport {
		transports = append(transports, string(t))
	}

	tx, err := server.DBConnection.Begin()
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer tx.Rollback()
	qtx := server.DBClient.WithTx(tx)

	user, err := qtx.UserAdd(ctx, db.UserAddParams{Username: userEntity.Username, ID: sessionData.UserID})
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	err = qtx.UserCreatePasskey(ctx, db.UserCreatePasskeyParams{
		PasskeyID:       credential.ID,
		UserID:          user.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Aaguid:          credential.Authenticator.AAGUID,
		SignCount:       int64(credential.Authenticator.SignCount),
		Transports:      transports,
		Name:            request.PasskeyName,
		BackupEligible:  credential.Flags.BackupEligible,
		BackupState:     credential.Flags.BackupState,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// if err := server.ApproveSession(ctx, user); err != nil {
	// 	log.Error(err)
	// 	return &webauthnproto.FinishResponse{}, err
	// }

	server.WebAuthnSessionStore.Delete(request.GetSessionId())

	return nil, nil
}

// BeginSignIn implements [v1connect.RakerServerHandler].
func (server *RakerServer) BeginSignIn(ctx context.Context, request *passkey.BeginSignInRequest) (*passkey.BeginSignInResponse, error) {
	username := request.GetUsername()
	var options *protocol.CredentialAssertion
	var sessionData *webauthn.SessionData
	var err error

	if username != "" {
		// Traditional non-discoverable flow with explicit username lookup
		user, err := server.DBClient.UserGetByUsername(ctx, username)
		if err != nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
		}
		userEntity := &authenticator.UserEntity{ID: user.ID, Username: user.Username}
		options, sessionData, err = server.WebAuthn.BeginLogin(userEntity)
	} else {
		// Discoverable/Passkey flow (allowCredentials will be empty)
		options, sessionData, err = server.WebAuthn.BeginDiscoverableLogin()
	}

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	sessionID := uuid.NewString()
	server.WebAuthnSessionStore.Set(sessionID, authenticator.WebAuthnSession{WebAuthnSession: *sessionData})

	optionsJSON, _ := json.Marshal(options.Response)

	return &passkey.BeginSignInResponse{
		SessionId:   sessionID,
		OptionsJson: string(optionsJSON),
	}, nil
}

// FinishSignIn implements [v1connect.RakerServerHandler].
func (server *RakerServer) FinishSignIn(ctx context.Context, request *passkey.FinishSignInRequest) (*emptypb.Empty, error) {
	sessionData, ok := server.WebAuthnSessionStore.Get(request.GetSessionId())
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("session expired or invalid"))
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "", bytes.NewReader([]byte(request.GetResponseJson())))
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	var loggedInUser *db.User

	credential, err := server.WebAuthn.FinishDiscoverableLogin(
		func(rawID []byte, userHandle []byte) (webauthn.User, error) {
			// Extract UUID from passkey's userHandle
			userID, err := uuid.FromBytes(userHandle)
			if err != nil {
				log.Error(err)
				return nil, fmt.Errorf("invalid userHandle: %w", err)
			}
			// log.Printf("Lookup User ID: %s", userID)

			// Look up user in database
			user, err := server.DBClient.UserGetByID(ctx, userID)
			if err != nil {
				log.Error(err)
				return nil, fmt.Errorf("user not found: %w", err)
			}
			loggedInUser = &user

			// Fetch passkeys belonging to this user
			dbKeys, err := server.DBClient.UserGetPasskeysByID(ctx, user.ID)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			// log.Printf("Found %d keys in DB for user", len(dbKeys))

			var credentials []webauthn.Credential
			for _, k := range dbKeys {
				credentials = append(credentials, webauthn.Credential{
					ID:              k.ID,
					PublicKey:       k.PublicKey,
					AttestationType: k.AttestationType,
					Authenticator:   webauthn.Authenticator{SignCount: uint32(k.SignCount)},
					Flags: webauthn.CredentialFlags{
						BackupEligible: k.BackupEligible,
						BackupState:    k.BackupState,
					},
				})
			}

			// Return user entity containing their public keys for signature verification
			return &authenticator.UserEntity{
				ID:       user.ID,
				Username: user.Username,
				Passkeys: credentials,
			}, nil
		},
		sessionData.WebAuthnSession,
		httpReq,
	)
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("login failed: %w", err))
	}

	if err := server.DBClient.PasskeyUpdateSignCount(ctx, credential.ID); err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := server.ApproveSession(ctx, *loggedInUser); err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	server.WebAuthnSessionStore.Delete(request.GetSessionId())

	return nil, nil
}

// RenamePasskey implements [v1connect.RakerServerHandler].
func (server *RakerServer) RenamePasskey(ctx context.Context, request *passkey.Passkey) (*emptypb.Empty, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	err := server.DBClient.PasskeyRename(ctx, db.PasskeyRenameParams{
		Name:      request.Name,
		PasskeyID: request.Id,
		UserID:    user.ID,
	})
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return nil, nil
}

// DeletePasskey implements [v1connect.RakerServerHandler].
func (server *RakerServer) DeletePasskey(ctx context.Context, request *passkey.Passkey) (*emptypb.Empty, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	passkeys, err := server.DBClient.UserGetPasskeysByID(ctx, user.ID)
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if len(passkeys) < 2 {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("you have only one passkey, you'll need to add another one before deleting this passkey."))
	}

	err = server.DBClient.PasskeyDelete(ctx, db.PasskeyDeleteParams{
		PasskeyID: request.Id,
		UserID:    user.ID,
	})
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return nil, nil
}

// GetPasskeysList implements [v1connect.RakerServerHandler].
func (server *RakerServer) GetPasskeysList(ctx context.Context, request *emptypb.Empty) (*passkey.PasskeysSettingsDisplay, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	passkeys, err := server.DBClient.UserGetPasskeysByID(ctx, user.ID)
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	passkeySettingsList := make([]*passkey.Passkey, 0, len(passkeys))
	for _, p := range passkeys {
		passkeySettingsList = append(passkeySettingsList, &passkey.Passkey{
			Id:     p.ID,
			Name:   p.Name,
			Aaguid: p.Aaguid,
		})
	}

	return &passkey.PasskeysSettingsDisplay{Passkeys: passkeySettingsList}, nil
}
