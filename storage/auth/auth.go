package auth

import srv "github.com/im-tollu/go-musthave-diploma-tpl/service/auth"

type Storage interface {
	CreateUser(u srv.UserToCreate) (srv.User, error)
	GetUserByLogin(login string) (*srv.User, error)
	SetUserSession(u srv.UserSessionToStart) (srv.UserSession, error)
	GetUserSession(uID int64) (srv.UserSession, error)
}
