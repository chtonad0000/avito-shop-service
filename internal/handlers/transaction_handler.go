package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/avito-shop-service/internal/models"
	"net/http"
	"time"

	"github.com/avito-shop-service/internal/middleware"
	"github.com/avito-shop-service/internal/services"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type TransactionHandler struct {
	userService        services.UserServiceInterface
	transactionService services.TransactionServiceInterface
}

func NewTransactionHandler(userService services.UserServiceInterface, transactionService services.TransactionServiceInterface) *TransactionHandler {
	return &TransactionHandler{
		userService:        userService,
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	fromUser, ok := middleware.GetEmployeeUsername(r.Context())
	if !ok {
		http.Error(w, "user not authorized", http.StatusUnauthorized)
		return
	}

	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.ToUser == "" || req.Amount <= 0 {
		http.Error(w, "invalid toUser or amount", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByUsername(r.Context(), fromUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	if user.Coins < req.Amount {
		http.Error(w, "not enough money to send", http.StatusBadRequest)
		return
	}

	err = h.transactionService.CreateTransaction(r.Context(), &models.CoinTransaction{
		UserID:          user.ID,
		CounterpartUser: req.ToUser,
		Amount:          req.Amount,
		TransactionType: "send",
		CreatedAt:       time.Now(),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to send coins: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "send successful"})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}
