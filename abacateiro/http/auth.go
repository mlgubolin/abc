package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterAuthRoutes(router chi.Router) {
	router.Post("/login", s.handleLogin)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Decode

	// Login user

	// Generate token

	// Send to Frontend
}
