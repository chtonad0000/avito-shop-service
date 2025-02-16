package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/avito-shop-service/internal/models"
	"net/http"

	"github.com/avito-shop-service/internal/services"
)

type UserHandler struct {
	service services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{service: userService}
}

func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}
	if user == nil {
		_, err = h.service.CreateUser(r.Context(), req.Username, req.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
			return
		}
	}
	token, err := h.service.Authenticate(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid username or password: %v", err), http.StatusUnauthorized)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}
