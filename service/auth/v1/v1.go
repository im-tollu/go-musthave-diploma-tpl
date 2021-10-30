package v1

import (
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/auth"
	authStorage "github.com/im-tollu/go-musthave-diploma-tpl/storage/auth"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

type Service struct {
	storage authStorage.Storage
}

func NewService() (*Service, error) {
	srv := Service{}

	return &srv, nil
}

func (s *Service) Register(cred auth.Credentials) error {
	hash, errHash := bcrypt.GenerateFromPassword(cred.Password, bcryptCost)
	if errHash != nil {
		return fmt.Errorf("cannot hash password: %w", errHash)
	}

	u := authStorage.UserToCreate{
		Login:        cred.Login,
		PasswordHash: hash,
	}

	if _, errCreate := s.storage.CreateUser(u); errCreate != nil {
		return fmt.Errorf("cannot create user [%s]: %w", cred.Login, errCreate)
	}

	return nil
}

func (s *Service) Login(cred auth.Credentials) (auth.SignedUserID, error) {
	nilLogin := auth.SignedUserID{}

	u, errGet := s.storage.GetUserByLogin(cred.Login)
	if errGet != nil {
		return nilLogin, fmt.Errorf("cannot get user by sess [%s]: %w", cred.Login, errGet)
	}
	if u == nil {
		return nilLogin, fmt.Errorf("user not found by sess [%s]: %w", cred.Login, auth.ErrWrongCredentials)
	}

	if errValidate := bcrypt.CompareHashAndPassword(u.PasswordHash, cred.Password); errValidate != nil {
		return nilLogin, fmt.Errorf("password doesn't match for user [%s]: %w", cred.Login, auth.ErrWrongCredentials)
	}

	sigKey, errGenKey := generateKey()
	if errGenKey != nil {
		return nilLogin, fmt.Errorf("cannot generate signature key: %w", errGenKey)
	}

	sess := authStorage.UserSession{
		UserID:       u.ID,
		SignatureKey: sigKey,
	}

	if errSetKey := s.storage.SetUserSession(sess); errSetKey != nil {
		return nilLogin, fmt.Errorf("cannot create session for user [%d]: %w", u.ID, errSetKey)
	}

	signedUserID := signUserId(sess)

	return signedUserID, nil
}
