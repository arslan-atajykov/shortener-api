package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/arslan-atajykov/shortener-api/internal/db"
	"github.com/go-chi/chi/v5"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_code = $1`

	err := db.DB.QueryRow(ctx, query, code).Scan(&originalURL)
	if err != nil {
		http.NotFound(w, r)
		return

	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
