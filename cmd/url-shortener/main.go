package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
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
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

	// todo: init router: chi, 'chi render'
	router := chi.NewRouter()
	// Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.RequestID)
	// Логирование всех запросов
	router.Use(middleware.Logger)
	//router.Use(mwLogger.New(log))
	// Если внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.Recoverer)
	// Парсер URLов поступающих запросов
	router.Use(middleware.URLFormat)

	router.Post("/", save.New(log, storage))

	router.Get("/{alias}", redirect.New(log, storage))

	//router.Delete("/{alias}", remove.New(log, storage))

	//-------------------------------

	// Все пути этого роутера будут начинаться с префикса `/url`
	router.Route("/url", func(r chi.Router) {
		// Подключаем авторизацию
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			// Передаем в middleware креды
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
			// Если у вас более одного пользователя,
			// то можете добавить остальные пары по аналогии.
		}))

		r.Post("/", save.New(log, storage))
	})

	// Хэндлер redirect остается снаружи, в основном роутере
	router.Get("/{alias}", redirect.New(log, storage))

	//-------------------------------

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
