package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const rememberTokenBytes = 128

func Bytes(n uint) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func String(n uint) (string, error) {
	bytes, err := Bytes(n)
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func RememberToken() (string, error) {
	return String(rememberTokenBytes)
}