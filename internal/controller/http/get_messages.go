package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// GetMessages implements the get_messages endpoint on the controller layer.
// The method parses request structure and calls usecase.Usecase.GetMessages().
func (c *HTTPController) GetMessages(w http.ResponseWriter, r *http.Request) {
	isIdentified := r.Context().Value("is_identified").(bool)
	if !isIdentified {
		http.Error(w, "not identified", http.StatusUnauthorized)
		return
	}

	userID := r.Context().Value("user_id").(uint)

	chatID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	var input dto.GetMessagesInput = dto.GetMessagesInput{UserID: userID, ChatID: uint(chatID)}

	if err := c.v.Struct(input); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	output, err := c.uc.GetMessages(input)
	if err != nil {
		if errors.Is(err, domain.ErrUserIsNotInChat) {
			http.Error(w, "you is not in the chat", http.StatusForbidden)
			return
		}
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
