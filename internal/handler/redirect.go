package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/arslan-atajykov/shortener-api/internal/cache"
	"github.com/arslan-atajykov/shortener-api/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	originalURL, err := cache.Client.Get(ctx, code).Result()
	if err == redis.Nil {
		var dbURL string
		query := `SELECT original_url FROM urls WHERE short_code = $1`
		err := db.DB.QueryRow(ctx, query, code).Scan(&dbURL)
		if err != nil {
			log.Printf("[REDIRECT] Code '%s' not found in DB", code)
			http.NotFound(w, r)
			return
		}

		// Кэшируем
		err = cache.Client.Set(ctx, code, dbURL, 24*time.Hour).Err()
		if err != nil {
			log.Printf("[REDIS] Failed to cache code '%s': %v", code, err)
		}

		originalURL = dbURL
		log.Printf("[REDIRECT] (from DB) %s → %s", code, originalURL)

	} else if err != nil {
		log.Printf("[REDIS] Error fetching code '%s': %v", code, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	} else {
		log.Printf("[REDIRECT] (from cache) %s → %s", code, originalURL)
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
