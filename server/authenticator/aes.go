package authenticator

import (
	"crypto/aes"
	"encoding/hex"
)

func Encrypt(key, pt string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	out := make([]byte, len(pt))
	c.Encrypt(out, []byte(pt))
	return hex.EncodeToString(out), nil
}

func Decrypt(key, ct string) (string, error) {
	in, err := hex.DecodeString(ct)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	pt := make([]byte, len(ct))
	c.Decrypt(pt, in)
	return string(pt[:]), nil
}
