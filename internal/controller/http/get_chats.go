package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"net/http"
)

// getChatsRequest defines signup request structure
type getChatsRequest = dto.GetChatsInput

// GetChats implements the get_chats endpoint on the controller layer.
// The method parses request structure and calls usecase.Usecase.GetChats().
func (c *HTTPController) GetChats(w http.ResponseWriter, r *http.Request) {
	isIdentified := r.Context().Value("is_identified").(bool)
	if !isIdentified {
		http.Error(w, "not identified", http.StatusUnauthorized)
		return
	}

	userID := r.Context().Value("user_id").(uint)

	var input dto.GetChatsInput = dto.GetChatsInput{UserID: userID}

	output, err := c.uc.GetChats(input)
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
