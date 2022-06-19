package hmac

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"hash"
)

// HMAC is a wrapper around crypto/hmac package making it a little easier to use
type HMAC struct {
	hmac hash.Hash
}

// NewHMAC creates and returns a new HMAC object
func NewHMAC(key string) HMAC {
	h := hmac.New(sha512.New, []byte(key))

	return HMAC{
		hmac: h,
	}
}

// Hash will hash the given input using HMAC algorithm
// and the secret key provided when the HMAC object was created
func (h HMAC) Hash(input string) (string, error) {
	h.hmac.Reset()
	_, err := h.hmac.Write([]byte(input)); if err != nil {
		return "", err
	}
	hash := h.hmac.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash), nil
}