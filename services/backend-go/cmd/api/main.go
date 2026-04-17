package main

import (
	"log/slog"
	"os"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/app"
)

func main() {
	cfg := app.LoadConfigFromEnv()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	server, err := app.NewServer(cfg, logger)
	if err != nil {
		logger.Error("server_init_error", "error", err)
		os.Exit(1)
	}
	if err := server.Run(); err != nil {
		logger.Error("server_exit", "error", err)
		os.Exit(1)
	}
}
