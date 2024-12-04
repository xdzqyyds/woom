package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"woom/server/model"

	"github.com/redis/go-redis/v9"
)

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 解析请求体中的JSON数据
	var inviteData struct {
		MeetingId string `json:"meetingId"`
		InviterId string `json:"inviterId"`
		InviteeId string `json:"inviteeId"`
	}

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
	response := map[string]string{
		"success": "true",
		"message": "Invitation sent successfully",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody struct {
		InviteeId string `json:"inviteeId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	inviteeId := reqBody.InviteeId
	value, err := h.rdb.HGet(ctx, model.InvitationKey, inviteeId).Result()

	if err == redis.Nil {
		return
	} else if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"value": value}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	if err := h.rdb.HDel(ctx, model.InvitationKey, inviteeId).Err(); err != nil {
		log.Printf("Failed to delete inviteeId %s from invitation: %v\n", inviteeId, err)
	}
}
