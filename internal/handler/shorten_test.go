package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestShortenHandler_ValidURL(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/shorten", ShortenHandler)

	body := map[string]string{
		"url": "https://example.com",
	}

	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp shortenResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.ShortURL, "http://localhost:8080/")
}
