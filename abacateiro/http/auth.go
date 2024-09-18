package http

import (
	"application"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//https://medium.com/@cheickzida/golang-implementing-jwt-token-authentication-bba9bfd84d60

func (s *Server) RegisterAuthRoutes(router chi.Router) {
	router.Post("/login", s.handleLogin)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Decode
	var loginUserQuery *application.LoginUserQuery

	if err := json.NewDecoder(r.Body).Decode(loginUserQuery); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "Invalid request payload"}) //Manda via http em vez de print
		return
	}

	// Login user
	userInfo, err := s.authService.Login(r.Context(), loginUserQuery)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()}) //Manda via http em vez de print
		return
	}

	// Generate token
	token, err := s.tokenService.GenerateToken(userInfo)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()}) //Manda via http em vez de print
		return
	}
	// Send to Frontend

	authInfo := map[string]interface{}{
		"token":   token.Token,
		"expiry":  token.Expiry,
		"email":   userInfo.Email,
		"auth_id": userInfo.AuthId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authInfo)
}
