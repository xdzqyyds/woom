package v1

import (
	"context"
	"encoding/json"
	"fmt"
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

type InviteRequest struct {
	MeetingId string `json:"meetingId"`
	InviterId string `json:"inviterId"`
	InviteeId string `json:"inviteeId"`
}

type InviteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetInvitationRequest struct {
	InviteeId string `json:"inviteeId"`
}

type UpdateUserStatusRequest struct {
	UserID string `json:"userId"`
	Status string `json:"status"`
}

type RemoveStreamRequest struct {
	StreamId string `json:"streamId"`
}

type MakeRemoveRequest struct {
	StreamId string `json:"streamId"`
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

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var inviteData InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&inviteData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	inviteValue := fmt.Sprintf("%s %s", inviteData.MeetingId, inviteData.InviterId)
	err := h.rdb.HSet(ctx, model.InvitationKey, inviteData.InviteeId, inviteValue).Err()
	if err != nil {
		log.Printf("Failed to save invitation for user %s: %v", inviteData.InviteeId, err)
		http.Error(w, "Failed to store invitation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := InviteResponse{
		Success: true,
		Message: "Invitation sent successfully",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody GetInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	value, err := h.rdb.HGet(ctx, model.InvitationKey, reqBody.InviteeId).Result()
	if err == redis.Nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"value": nil,
		})
		return
	} else if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"value": value}); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	if err := h.rdb.HDel(ctx, model.InvitationKey, reqBody.InviteeId).Err(); err != nil {
		log.Printf("Failed to delete inviteeId %s from invitation: %v\n", reqBody.InviteeId, err)
	}
}

func (h *Handler) UserList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	onlineStatus, err := h.rdb.HGetAll(ctx, model.UserOnlineStatusKey).Result()
	if err != nil {
		http.Error(w, "Failed to get user online status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(onlineStatus); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateUserList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestData UpdateUserStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	err := h.rdb.HSet(ctx, model.UserOnlineStatusKey, requestData.UserID, requestData.Status).Err()
	if err != nil {
		http.Error(w, "Failed to update user status", http.StatusInternalServerError)
		return
	}
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

func (h *Handler) RemoveStream(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req RemoveStreamRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.rdb.SAdd(ctx, model.StreamRemovalKey, req.StreamId).Err(); err != nil {
		http.Error(w, "Failed to add stream to removal set", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) MakeRemove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req MakeRemoveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	isMember, err := h.rdb.SIsMember(ctx, model.StreamRemovalKey, req.StreamId).Result()
	if err != nil {
		http.Error(w, "Failed to check stream removal set", http.StatusInternalServerError)
		return
	}

	value := 0
	if isMember {
		value = 1
		if err := h.rdb.SRem(ctx, model.StreamRemovalKey, req.StreamId).Err(); err != nil {
			http.Error(w, "Failed to remove stream from set", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"value": value})
}
