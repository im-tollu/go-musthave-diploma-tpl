package handler

import (
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/apimodel"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/auth"
	"net/http"
)

func makeAuthCookie(u auth.SignedUserID) http.Cookie {
	v := fmt.Sprintf("%d|%x", u.ID, u.Signature)
	return http.Cookie{
		Name:  apimodel.AuthCookieName,
		Value: v,
	}
}
