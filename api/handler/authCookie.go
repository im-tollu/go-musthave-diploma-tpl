package handler

import (
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/apiModel"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/auth"
	"net/http"
)

func makeAuthCookie(u auth.SignedUserID) http.Cookie {
	v := fmt.Sprintf("%d|%x", u.ID, u.Signature)
	return http.Cookie{
		Name:  apiModel.AuthCookieName,
		Value: v,
	}
}
