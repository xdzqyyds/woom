package v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Invalid request"})
		return
	}

	// Redis hash key
	userStorageKey := "user_storage"

	// get password
	ctx := context.Background()
	passwordKey := loginReq.Username
	log.Printf("Attempting to fetch password for key: %s\n", passwordKey)
	password, err := h.rdb.HGet(ctx, userStorageKey, passwordKey).Result()
	if err != nil {
		log.Printf("Failed to find password for key: %s, error: %v\n", passwordKey, err)
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, LoginResponse{Success: false, Message: "User not found"})
		return
	}

	// verify password
	if password != loginReq.Password {
		log.Printf("Password mismatch: expected %s, got %s\n", password, loginReq.Password)
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Invalid password"})
		return
	}

	render.JSON(w, r, LoginResponse{Success: true, Message: "Login successful"})
}
