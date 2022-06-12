package authenticator

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

func (a *Authenticator) Encrypt(pt string) (string, error) {
	sum := sha256.Sum256([]byte(a.secret))
	c, err := aes.NewCipher(sum[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ct := hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(pt), nil))
	return ct, nil
}

func (a *Authenticator) Decrypt(ct string) (string, error) {
	sum := sha256.Sum256([]byte(a.secret))
	c, err := aes.NewCipher(sum[:])
	if err != nil {
		return "", err
	}

	ctBytes, err := hex.DecodeString(ct)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	if len(ctBytes) < gcm.NonceSize() {
		return "", errors.New("cipher text too short")
	}

	nonce, ctBytes := ctBytes[:gcm.NonceSize()], ctBytes[gcm.NonceSize():]
	ptBytes, err := gcm.Open(nil, nonce, ctBytes, nil)
	return string(ptBytes), err
}
