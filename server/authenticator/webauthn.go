package authenticator

import (
	"context"
	"sync"
	"time"

	"github.com/AppleGamer22/raker/shared"
	"github.com/charmbracelet/log"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
)

type WebAuthnSession struct {
	WebAuthnSession webauthn.SessionData
	UserID          uuid.UUID
	Username        string
}

type WebAuthnSessionStore struct {
	mu      sync.Mutex
	cache   *ttlcache.Cache[string, WebAuthnSession]
	running bool
}

func NewSessionStore() *WebAuthnSessionStore {
	store := WebAuthnSessionStore{}

	store.cache = ttlcache.New(
		ttlcache.WithTTL[string, WebAuthnSession](5 * time.Minute),
	)

	// Stop cache cleaner when eviction empties active sessions
	store.cache.OnEviction(func(ctx context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[string, WebAuthnSession]) {
		store.mu.Lock()
		defer store.mu.Unlock()

		if store.cache.Len() == 0 && store.running {
			store.running = false
			store.cache.Stop()
			log.Debug("stopping webauthn session cache")
		}
	})

	return &store
}

func (s *WebAuthnSessionStore) Set(key string, session WebAuthnSession) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Lazy start if cache is not running
	if !s.running {
		s.running = true
		go s.cache.Start()
		log.Debug("starting webauthn session cache")
	}
	if shared.Version == "development" {
		log.Debug("set session", "id", key)
	}
	s.cache.Set(key, session, ttlcache.DefaultTTL)
}

func (s *WebAuthnSessionStore) Get(key string) (WebAuthnSession, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := s.cache.Get(key)
	if item == nil {
		return WebAuthnSession{}, false
	}

	data := item.Value()

	return data, true
}

func (s *WebAuthnSessionStore) Delete(key string) {
	if shared.Version == "development" {
		log.Debug("delete session", "id", key)
	}
	s.cache.Delete(key)
}

type UserEntity struct {
	ID       uuid.UUID
	Username string
	Passkeys []webauthn.Credential
}

func (u *UserEntity) WebAuthnID() []byte {
	return u.ID[:]
}

func (u *UserEntity) WebAuthnName() string {
	return u.Username
}

func (u *UserEntity) WebAuthnDisplayName() string {
	return u.Username
}

func (u *UserEntity) WebAuthnIcon() string {
	return ""
}

func (u *UserEntity) WebAuthnCredentials() []webauthn.Credential {
	return u.Passkeys
}
