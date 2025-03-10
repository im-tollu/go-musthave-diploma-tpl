package api

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/handler"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/middleware"
)

type loyaltyRouter struct {
	*chi.Mux
}

func newRouter(h *handler.LoyaltyHandler) *loyaltyRouter {
	router := loyaltyRouter{
		Mux: chi.NewMux(),
	}

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Group(func(g chi.Router) {
			g.Use(chimiddleware.Logger)
			g.Use(middleware.Authenticator(h.AuthSrv))
			g.Post("/orders", h.PostOrder)
			g.Get("/orders", h.GetOrders)
			g.Get("/balance", h.Balance)
			g.Post("/balance/withdraw", h.Withdraw)
			g.Get("/withdrawals", h.GetWithdrawals)
		})
	})

	return &router
}
