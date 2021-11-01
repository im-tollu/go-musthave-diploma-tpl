package api

import (
	"context"
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/handler"
	auth "github.com/im-tollu/go-musthave-diploma-tpl/service/auth/v1"
	order "github.com/im-tollu/go-musthave-diploma-tpl/service/order/v1"
	authStorage "github.com/im-tollu/go-musthave-diploma-tpl/storage/auth"
	orderStorage "github.com/im-tollu/go-musthave-diploma-tpl/storage/order"
	"log"
	"net/http"
)

type LoyaltyServer struct {
	http.Server
}

// NewServer makes an instance of LoyaltyServer HTTP server and runs it
// in a separate goroutine
func NewServer(addr string, authStorage authStorage.Storage, orderStorage orderStorage.Storage) (*LoyaltyServer, error) {
	authSrv, errAuth := auth.NewService(authStorage)
	if errAuth != nil {
		return nil, fmt.Errorf("cannot get instance of Auth Service: %w", errAuth)
	}

	orderSrv, errOrder := order.NewService(orderStorage)
	if errOrder != nil {
		return nil, fmt.Errorf("cannot get instance of Order Service: %w", errOrder)
	}

	h, errHandler := handler.NewHandler(authSrv, orderSrv)
	if errHandler != nil {
		return nil, fmt.Errorf("cannot get instance of Handler: %w", errHandler)
	}

	r := newRouter(h)
	server := LoyaltyServer{
		Server: http.Server{
			Addr:    addr,
			Handler: r,
		},
	}

	log.Println("Starting server...")

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("Server failed: %s", err.Error())
		}
	}()

	return &server, nil
}

// Shutdown gracefully stops the server
func (s *LoyaltyServer) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server stopped.")

	return nil
}
