package handler

import (
	"github.com/im-tollu/go-musthave-diploma-tpl/api/middleware"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/auth"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/order"
	"net/http"
)

type LoyaltyHandler struct {
	AuthSrv  auth.Service
	OrderSrv order.Service
}

func NewHandler(authSrv auth.Service, orderSrv order.Service) (*LoyaltyHandler, error) {
	handler := LoyaltyHandler{
		AuthSrv:  authSrv,
		OrderSrv: orderSrv,
	}

	return &handler, nil
}

func userID(r *http.Request) int64 {
	return r.Context().Value(middleware.AuthContextKeyType{}).(int64)
}
