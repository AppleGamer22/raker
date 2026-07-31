package authenticator

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
)

type WebAuthnSessionStore struct {
	mu      sync.Mutex
	cache   *ttlcache.Cache[string, webauthn.SessionData]
	running bool
}

func NewSessionStore() WebAuthnSessionStore {
	store := WebAuthnSessionStore{}

	store.cache = ttlcache.New(
		ttlcache.WithTTL[string, webauthn.SessionData](5 * time.Minute),
	)

	// Stop cache cleaner when eviction empties active sessions
	store.cache.OnEviction(func(ctx context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[string, webauthn.SessionData]) {
		store.mu.Lock()
		defer store.mu.Unlock()

		if store.cache.Len() == 0 && store.running {
			store.running = false
			store.cache.Stop()
			log.Debug("stopping webauthn session cache")
		}
	})

	return store
}

func (s *WebAuthnSessionStore) Set(key string, session webauthn.SessionData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Lazy start if cache is not running
	if !s.running {
		s.running = true
		go s.cache.Start()
		log.Debug("starting webauthn session cache")
	}

	s.cache.Set(key, session, ttlcache.DefaultTTL)
}

func (s *WebAuthnSessionStore) GetAndDelete(key string) (webauthn.SessionData, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := s.cache.Get(key)
	if item == nil {
		return webauthn.SessionData{}, false
	}

	data := item.Value()
	s.cache.Delete(key)

	return data, true
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
