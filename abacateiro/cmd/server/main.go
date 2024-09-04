package main

import (
    "application/http"
    "application/postgres"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func dd(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Erro ao serializar dados:", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
	os.Exit(0)
}


func initializeLogger() *log.Logger {
    return log.New(log.Writer(), "app: ", log.LstdFlags)
}

func initializeServer(logger *log.Logger, userService *postgres.UserService) *http.Server {
    return http.NewServer(":8080", logger, userService)
}

func main() {
	
    fmt.Println("Iniciando o servidor...")

    logger := initializeLogger()
    userService := postgres.NewUserService()
    server := initializeServer(logger, userService)

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
