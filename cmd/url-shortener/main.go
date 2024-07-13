package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
)

func main() {
	// todo: init config: cleanenv
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// todo: init logger: slog
	log := setupLogger(cfg.Env)
	log.Info("start url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// todo: init storage: sqlite / postgres

	// todo: init router: chi, 'chi render'

	// todo: run server
}

const (
	envLocal = "local"
	envDel   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDel:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
