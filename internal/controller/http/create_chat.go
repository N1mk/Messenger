package http

import (
	"encoding/json"
	"fmt"
	"messenger/internal/dto"
	"net/http"
)

// createChatRequest defines signup request structure
type createChatRequest struct {
	Usernames []string `json:"usernames" validate:"required"`
	Name      string   `json:"name" validate:"required,max=20"`
}

// CreateChat implements the create_chat endpoint on the controller layer.
// The method parses request structure and calls usecase.Usecase.CreateChat().
func (c *HTTPController) CreateChat(w http.ResponseWriter, r *http.Request) {
	isIdentified := r.Context().Value("is_identified").(bool)
	if !isIdentified {
		http.Error(w, "not identified", http.StatusUnauthorized)
		return
	}

	creatorID := r.Context().Value("user_id").(uint)

	var req createChatRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("decoding error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var input dto.CreateChatInput = dto.CreateChatInput{CreatorID: creatorID, Usernames: req.Usernames, Name: req.Name}

	output, err := c.uc.CreateChat(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, fmt.Sprintf("encoding error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
