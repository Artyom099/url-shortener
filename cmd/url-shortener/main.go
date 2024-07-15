package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	//router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer) // поднимает прилодение, если сработал panic
	router.Use(middleware.URLFormat)

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
		//log = setupPrettySlog()
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDel:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
