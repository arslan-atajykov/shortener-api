package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arslan-atajykov/shortener-api/internal/config"
	"github.com/arslan-atajykov/shortener-api/internal/db"
	"github.com/arslan-atajykov/shortener-api/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	db.Init(cfg)

	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})
	r.Post("/shorten", handler.ShortenHandler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Println("Server running on", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed:", err)
	}

}
