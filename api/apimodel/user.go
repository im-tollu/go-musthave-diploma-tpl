package apimodel

import "github.com/im-tollu/go-musthave-diploma-tpl/service/auth"

const AuthCookieName = "USER-AUTH"

type CredentialsJSON struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (cr CredentialsJSON) ToCredentials() auth.Credentials {
	return auth.Credentials{
		Login:    cr.Login,
		Password: []byte(cr.Password),
	}
}
