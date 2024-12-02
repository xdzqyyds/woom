package v1

import (
	"encoding/json"
	"net/http"
	"woom/server/model"
)

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

	var requestData struct {
		UserID string `json:"userId"`
		Status string `json:"status"`
	}

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
