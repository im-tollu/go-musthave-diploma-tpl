package handler

import "github.com/im-tollu/go-musthave-diploma-tpl/service/auth"

type LoyaltyHandler struct {
	authSrv auth.Service
}

func NewHandler(authSrv auth.Service) (*LoyaltyHandler, error) {
	handler := LoyaltyHandler{
		authSrv: authSrv,
	}

	return &handler, nil
}
