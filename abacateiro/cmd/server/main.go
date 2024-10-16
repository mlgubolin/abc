package main

import (
	"application/auth"
	"application/http"
	"application/postgres"
	"application/token"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// func dd(data interface{}) {
// 	jsonData, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		fmt.Println("Erro ao serializar dados:", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(jsonData))
// 	os.Exit(0)
// }

func initializeLogger() *log.Logger {
	return log.New(log.Writer(), "app: ", log.LstdFlags)
}

func initializeServer(
	logger *log.Logger, userService *postgres.UserService,
	authService *auth.AuthService, tokenService *token.TokenService, workReportService *postgres.WorkReportService,
) *http.Server {

	return http.NewServer(":8080", logger, userService, authService, tokenService, workReportService)

}

func initializeDatabase(logger *log.Logger) *pgxpool.Pool {
	connString := "host=postgres port=5432 user=abacateiro password=abacateiro dbname=abacateiro sslmode=disable"
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		logger.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	return dbPool
}

func main() {
	fmt.Println("Iniciando o servidor...")

	logger := initializeLogger()
	dbPool := initializeDatabase(logger)
	userService := postgres.NewUserService(dbPool)
	authService := auth.NewAuthService(userService)
	tokenService := token.NewTokenService()
	workReportService := postgres.NewWorkReportService(dbPool)

	server := initializeServer(logger, userService, authService, tokenService, workReportService)

	// _, err := work_report.UnzipWorkReport("word.doc")
	// if err != nil {
	// 	return
	// }
	// Usar contexto para gerenciar o ciclo de vida do servidor
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.Start(); err != nil {
			logger.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Println("Desligando o servidor...")

	// Tempo de espera para o servidor desligar
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		logger.Fatalf("Erro ao desligar o servidor: %v", err)
	}

	logger.Println("Servidor desligado com sucesso")
}
