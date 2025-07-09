package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/arslan-atajykov/shortener-api/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (username, password_hash)
				VALUES ($1, $2)
				RETURNING id;`
	var id int
	err = db.DB.QueryRow(ctx, query, req.Username, string(hashed)).Scan(&id)
	if err != nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

var jwtSecret = []byte("mysecretkey")

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var storedHash string

	query := `SELECT password_hash FROM users WHERE username = $1`
	err := db.DB.QueryRow(ctx, query, req.Username).Scan(&storedHash)
	if err != nil {
		http.Error(w, "invalid username/password", http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid username/password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "couldnt generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(loginResponse{Token: tokenString})

}
