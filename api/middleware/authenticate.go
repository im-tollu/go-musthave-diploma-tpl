package middleware

import (
	"context"
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/apiModel"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/auth"
	"log"
	"net/http"
)

type AuthContextKeyType struct{}

// Authenticator middleware authenticates a request
// based on the signed cookie containing a user ID.
// In case authentication has failed, it signs up a new user.
func Authenticator(s auth.Service) func(http.Handler) http.Handler {
	ra := requestAuth{s}

	return func(next http.Handler) http.Handler {
		serveHTTP := func(w http.ResponseWriter, r *http.Request) {
			userID := ra.extractUserID(r)
			if userID == nil {
				http.Error(w, "Login to access this endpoint", http.StatusUnauthorized)
				return
			}

			ctxWithUserID := context.WithValue(r.Context(), AuthContextKeyType{}, *userID)

			next.ServeHTTP(w, r.WithContext(ctxWithUserID))
		}

		return http.HandlerFunc(serveHTTP)
	}
}

type requestAuth struct {
	AuthService auth.Service
}

func (a *requestAuth) extractUserID(r *http.Request) *int64 {
	cookie, errGetCookie := r.Cookie(apiModel.AuthCookieName)
	if errGetCookie != nil {
		log.Printf("Cannot get authentication cookie: %s", errGetCookie.Error())
		return nil
	}

	var userID int64
	var signature []byte
	if _, err := fmt.Sscanf(cookie.Value, "%d|%x", &userID, &signature); err != nil {
		log.Printf("Cannot parse authentication cookie [%s]: %s", cookie.Value, err.Error())
	}

	sgn := auth.SignedUserID{
		ID:        userID,
		Signature: signature,
	}

	if invalid := a.AuthService.Validate(sgn); invalid != nil {
		log.Printf("Signature is invalid: %s", invalid.Error())
		return nil
	}

	return &sgn.ID
}
