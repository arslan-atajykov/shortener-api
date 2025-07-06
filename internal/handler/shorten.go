package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/arslan-atajykov/shortener-api/internal/db"
	"github.com/arslan-atajykov/shortener-api/internal/service"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidURL(req.URL) {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var existingCode string
	checkQuery := `SELECT short_code FROM urls WHERE original_url = $1`
	//
	err := db.DB.QueryRow(ctx, checkQuery, req.URL).Scan(&existingCode)
	if err == nil {
		resp := shortenResponse{
			ShortURL: "http://localhost:8080/" + existingCode,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
	shortCode := service.GenerateShortCode()
	query := `INSERT INTO urls (original_url, short_code)
				VALUES($1, $2)
				RETURNING id;`

	var id int
	err = db.DB.QueryRow(ctx, query, req.URL, shortCode).Scan(&id)

	if err != nil {
		http.Error(w, "failed to save url", http.StatusInternalServerError)
		return
	}

	resp := shortenResponse{
		ShortURL: "http://localhost:8080/" + shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func isValidURL(raw string) bool {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
