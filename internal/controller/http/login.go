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

// loginRequest defines signup request structure
type loginRequest struct {
	Password string `json:"password" validate:"required,min=6,max=20"`
}

// Login implements the login endpoint on the controller layer.
// The method parses request structure and calls usecase.Usecase.Login().
func (c *HTTPController) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("decoding error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	var input dto.LoginInput = dto.LoginInput{UserID: uint(userID), Password: req.Password}

	if err := c.v.Struct(input); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	output, err := c.uc.Login(input)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			http.Error(w, "incorrect password", http.StatusBadRequest)
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
