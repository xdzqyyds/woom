package v1

import (
	"encoding/json"
	"net/http"
	"woom/server/model"
)

func (h *Handler) GetUserOnlineStatus(w http.ResponseWriter, r *http.Request) {
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
