package v1

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/auth"
	"strconv"
)

const KeySize = 16

func generateKey() ([]byte, error) {
	key := make([]byte, KeySize)
	if _, err := rand.Read(key); err != nil {
		return key, fmt.Errorf("cannot generate signature key: %w", err)
	}

	return key, nil
}

func signUserId(key auth.UserSession) auth.SignedUserID {
	h := hmac.New(sha256.New, key.SignatureKey)
	h.Write([]byte(strconv.FormatInt(key.UserID, 10)))
	sum := h.Sum(nil)

	return auth.SignedUserID{
		ID:        key.UserID,
		Signature: sum,
	}
}
