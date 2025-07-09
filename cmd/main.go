package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arslan-atajykov/shortener-api/internal/cache"
	"github.com/arslan-atajykov/shortener-api/internal/config"
	"github.com/arslan-atajykov/shortener-api/internal/db"
	"github.com/arslan-atajykov/shortener-api/internal/handler"
	"github.com/arslan-atajykov/shortener-api/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Загружаем конфигурации
	cfg := config.LoadConfig()

	// Инициализируем БД и кэш
	db.Init(cfg)
	cache.Init(cfg)

	// Создаем маршрутизатор
	r := chi.NewRouter()

	// Простой health check
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	// Аутентификация
	r.Post("/register", handler.RegisterHandler)
	r.Post("/login", handler.LoginHandler)

	// Публичный редирект
	r.Get("/{code}", handler.RedirectHandler)

	// Группа маршрутов с защитой JWT
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) // <-- JWT проверка
		r.Post("/shorten", handler.ShortenHandler)
	})

	// Запускаем сервер
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Println("Server running on", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed:", err)
	}
}
