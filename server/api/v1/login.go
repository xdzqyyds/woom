package v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"

	"woom/server/model"
)

type LoginRequest struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Invalid request"})
		return
	}

	// get password
	ctx := context.Background()
	passwordKey := loginReq.UserId
	password, err := h.rdb.HGet(ctx, model.UserStorageKey, passwordKey).Result()
	if err != nil {
		log.Printf("Failed to find password for key: %s, error: %v\n", passwordKey, err)
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, LoginResponse{Success: false, Message: "User not found"})
		return
	}

	// verify password
	if password != loginReq.Password {
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Invalid password"})
		return
	}

	// Update user online status to true
	if err := h.rdb.HSet(ctx, model.UserOnlineStatusKey, passwordKey, true).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Failed to update online status"})
		return
	}

	render.JSON(w, r, LoginResponse{Success: true, Message: "Login successful"})
}
