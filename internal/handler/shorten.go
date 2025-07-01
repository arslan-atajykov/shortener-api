// package handler

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"time"

// 	"github.com/arslan-atajykov/shortener-api/internal/db"
// )

// type shortenRequest struct {
// 	URL string `json:"url"`
// }

// type shortenResponse struct {
// 	ShortURL string `json:"short_url"`
// }

// func ShortenHandler(w http.ResponseWriter, r *http.Request) {
// 	var req shortenRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	shortCode := service.GenerateShortCode()
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	query := `INSERT INTO urls (original_url, short_code)
// 				VALUES($1, $2)
// 				RETURNING id;`

// 	var id int
// 	err := db.DB.QueryRow()

// }
