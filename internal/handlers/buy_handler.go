package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avito-shop-service/internal/middleware"
	"github.com/avito-shop-service/internal/services"
	"github.com/go-chi/chi/v5"
)

type BuyHandler struct {
	userService        services.UserServiceInterface
	merchService       services.MerchServiceInterface
	inventoryService   services.InventoryServiceInterface
	transactionService services.TransactionServiceInterface
}

func NewBuyHandler(userService services.UserServiceInterface, merchService services.MerchServiceInterface, inventoryService services.InventoryServiceInterface, transactionService services.TransactionServiceInterface) *BuyHandler {
	return &BuyHandler{
		userService:        userService,
		merchService:       merchService,
		inventoryService:   inventoryService,
		transactionService: transactionService,
	}
}

func (h *BuyHandler) Buy(w http.ResponseWriter, r *http.Request) {
	employeeUsername, ok := middleware.GetEmployeeUsername(r.Context())
	if !ok {
		http.Error(w, "user not authorized", http.StatusUnauthorized)
		return
	}
	itemName := chi.URLParam(r, "item")
	if itemName == "" {
		http.Error(w, "item name is required", http.StatusBadRequest)
		return
	}

	merch, err := h.merchService.GetMerchByName(r.Context(), itemName)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching merch: %v", err), http.StatusInternalServerError)
		return
	}
	if merch == nil {

		http.Error(w, "item not found", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByUsername(r.Context(), employeeUsername)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	if user.Coins < merch.Price {
		http.Error(w, "not enough coins", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.BuyItemToInventory(r.Context(), user.ID, merch.ID, 1, merch.Price)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating inventory: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "purchase successful"})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}
