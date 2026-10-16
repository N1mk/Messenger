package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"net/http"
)

// signupRequest defines signup request structure
type signupRequest = dto.SignupInput

// Signup implements the signup endpoint on the controller layer.
// The method parses request structure and calls usecase.Usecase.Signup().
func (c *HTTPController) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("decoding error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var input dto.SignupInput = req

	if err := c.v.Struct(input); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	output, err := c.uc.Signup(input)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicated) {
			http.Error(w, "duplicated", http.StatusBadRequest)
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

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, fmt.Sprintf("encoding error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
