package api

import (
	"context"
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/handler"
	auth "github.com/im-tollu/go-musthave-diploma-tpl/service/auth/v1"
	"log"
	"net/http"
)

type LoyaltyServer struct {
	http.Server
}

// NewServer makes an instance of LoyaltyServer HTTP server and runs it
// in a separate goroutine
func NewServer(addr string) (*LoyaltyServer, error) {
	authSrv, errAuth := auth.NewService()
	if errAuth != nil {
		return nil, fmt.Errorf("cannot get instance of Auth Service: %w", errAuth)
	}

	h, errHandler := handler.NewHandler(authSrv)
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
