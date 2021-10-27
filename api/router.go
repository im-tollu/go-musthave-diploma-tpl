package api

import "github.com/go-chi/chi/v5"

type loyaltyRouter struct {
	*chi.Mux
}

func newRouter() *loyaltyRouter {
	router := loyaltyRouter{
		Mux: chi.NewMux(),
	}

	return &router
}
