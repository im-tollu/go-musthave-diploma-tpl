package auth

import "errors"

var ErrLoginAlreadyTaken = errors.New("login already taken")

var ErrWrongCredentials = errors.New("incorrect login/password")

type Credentials struct {
	Login    string
	Password []byte
}

type SignedUserID struct {
	ID        int64
	Signature []byte
}
