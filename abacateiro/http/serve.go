package http

import (
	"application"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server estrutura principal do servidor HTTP
type Server struct {
	server       *http.Server
	logger       *log.Logger
	userService  application.UserService
	authService  application.AuthService
	tokenService application.TokenService
}

// NewServer construtor que inicializa um novo servidor HTTP
func NewServer(
	addr string,
	logger *log.Logger,
	userService application.UserService,
	authService application.AuthService,
	tokenService application.TokenService,
) *Server {

	// Criar o roteador chi
	router := chi.NewRouter()

	// Basic CORS settings
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:9000"}, // Or "*"
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Adicionar middlewares do chi
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Criar o servidor HTTP com o roteador chi
	serve := &Server{
		server: &http.Server{
			Addr:         addr,   // e.g. ":8080"
			Handler:      router, // roteador chi
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		logger:       logger,
		userService:  userService,
		authService:  authService,
		tokenService: tokenService,
	}

	// Adicionar rotas
	serve.registerRoutes(router)

	return serve
}

// registerRoutes adiciona as rotas ao roteador
func (s *Server) registerRoutes(router *chi.Mux) {
	s.RegisterUserRoutes(router)
	s.RegisterAuthRoutes(router)
}

// Start inicia o servidor HTTP
func (s *Server) Start() error {
	s.logger.Printf("Iniciando o servidor na porta %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatalf("Erro ao iniciar o servidor: %v", err)
		return err
	}
	return nil
}

// Stop finaliza o servidor HTTP
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Println("Parando o servidor")
	return s.server.Shutdown(ctx)
}
