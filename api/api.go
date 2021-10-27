package api

import (
	"context"
	"log"
	"net/http"
)

type LoyaltyServer struct {
	http.Server
}

// NewServer makes an instance of LoyaltyServer HTTP server and runs it
// in a separate goroutine
func NewServer(addr string) *LoyaltyServer {
	server := LoyaltyServer{
		Server: http.Server{
			Addr:    addr,
			Handler: newRouter(),
		},
	}

	log.Println("Starting server...")

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("Server failed: %s", err.Error())
		}
	}()

	return &server
}

// Shutdown gracefully stops the server
func (s *LoyaltyServer) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server stopped.")

	return nil
}
