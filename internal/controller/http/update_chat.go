package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"net/http"
)

// updateAccountRequest defines account update request structure
type updateChatRequest struct {
	ChatID    uint     `json:"chat_id" validate:"required"`
	CreatorID uint     `json:"user_id" validate:"required"`
	Usernames []string `json:"usernames" validate:"required"`
	Name      string   `json:"name" validate:"required,max=20"`
}

// UpdateAccount implements the update_account endpoint on the controller layer.
// The method parses request structure and calls usecase.Usecase.UpdateAccount().
func (c *HTTPController) UpdateChat(w http.ResponseWriter, r *http.Request) {
	isIdentified := r.Context().Value("is_identified").(bool)
	if !isIdentified {
		http.Error(w, "not identified", http.StatusUnauthorized)
		return
	}

	var req updateChatRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("decoding error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID := r.Context().Value("user_id").(uint)

	var input dto.UpdateChatInput = dto.UpdateChatInput{UserID: userID, ChatID: req.ChatID, CreatorID: req.CreatorID, Usernames: req.Usernames, Name: req.Name}

	if err := c.v.Struct(input); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := c.v.Struct(input); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	output, err := c.uc.UpdateChat(input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, fmt.Sprintf("encoding error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
