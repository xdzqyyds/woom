package v1

import (
	"context"
	"encoding/json"

	// "fmt"
	"log"
	"net/http"
	"woom/server/model"

	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
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

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Invalid request"})
		return
	}

	ctx := context.Background()
	passwordKey := loginReq.UserId

	onlineStatus, err := h.rdb.HGet(ctx, model.UserOnlineStatusKey, passwordKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Failed to get online status for user: %s, error: %v\n", passwordKey, err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Failed to check online status"})
		return
	}

	if onlineStatus == "1" {
		w.WriteHeader(http.StatusConflict)
		render.JSON(w, r, LoginResponse{Success: false, Message: "User already logged in"})
		return
	}

	password, err := h.rdb.HGet(ctx, model.UserStorageKey, passwordKey).Result()
	if err != nil {
		log.Printf("Failed to find password for key: %s, error: %v\n", passwordKey, err)
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, LoginResponse{Success: false, Message: "User not found"})
		return
	}

	if password != loginReq.Password {
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Invalid password"})
		return
	}

	if err := h.rdb.HSet(ctx, model.UserOnlineStatusKey, passwordKey, "1").Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Failed to update online status"})
		return
	}

	render.JSON(w, r, LoginResponse{Success: true, Message: "Login successful"})
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var signupReq LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&signupReq); err != nil {
		log.Printf("Error decoding signup request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Invalid request format"})
		return
	}

	// User input cannot be empty
	if signupReq.UserId == "" || signupReq.Password == "" {
		render.JSON(w, r, LoginResponse{Success: false, Message: "User ID and password cannot be empty"})
		return
	}

	ctx := context.Background()

	// Check whether the user name already exists
	exists, err := h.rdb.HExists(ctx, model.UserStorageKey, signupReq.UserId).Result()
	if err != nil {
		log.Printf("Error checking user existence: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Internal server error"})
		return
	}

	if exists {
		render.JSON(w, r, LoginResponse{Success: false, Message: "This user ID is already registered"})
		return
	}

	// Execute multiple Redis commands using pipes
	pipe := h.rdb.Pipeline()

	pipe.HSet(ctx, model.UserStorageKey, signupReq.UserId, signupReq.Password)
	pipe.HSet(ctx, model.UserOnlineStatusKey, signupReq.UserId, "0")

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Error executing pipeline for user registration: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, LoginResponse{Success: false, Message: "Failed to register user"})
		return
	}

	log.Printf("New user registered successfully: %s\n", signupReq.UserId)
	render.JSON(w, r, LoginResponse{Success: true, Message: "Registration successful! You can now login."})
}
