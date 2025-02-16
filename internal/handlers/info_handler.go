package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/avito-shop-service/internal/middleware"
	"github.com/avito-shop-service/internal/services"
	"net/http"
)

type InformationHandler struct {
	userService        services.UserServiceInterface
	merchService       services.MerchServiceInterface
	inventoryService   services.InventoryServiceInterface
	transactionService services.TransactionServiceInterface
}

func NewInformationHandler(userService services.UserServiceInterface, merchService services.MerchServiceInterface, inventoryService services.InventoryServiceInterface, transactionService services.TransactionServiceInterface) *InformationHandler {
	return &InformationHandler{
		userService:        userService,
		transactionService: transactionService,
		inventoryService:   inventoryService,
		merchService:       merchService,
	}
}

func (h *InformationHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	employeeUsername, ok := middleware.GetEmployeeUsername(r.Context())
	if !ok {
		http.Error(w, "user not authorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByUsername(r.Context(), employeeUsername)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, fmt.Sprintf("Error finding user: %v", err), http.StatusBadRequest)
		return
	}

	transactions, err := h.transactionService.GetTransactionsByUserId(r.Context(), user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching transactions: %v", err), http.StatusInternalServerError)
		return
	}

	inventory, err := h.inventoryService.GetInventoryByUserID(r.Context(), user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching inventory: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Coins     int `json:"coins"`
		Inventory []struct {
			Type     string `json:"type"`
			Quantity int    `json:"quantity"`
		} `json:"inventory"`
		CoinHistory struct {
			Received []struct {
				FromUser string `json:"fromUser"`
				Amount   int    `json:"amount"`
			} `json:"received"`
			Sent []struct {
				ToUser string `json:"toUser"`
				Amount int    `json:"amount"`
			} `json:"sent"`
		} `json:"coinHistory"`
	}{
		Coins: user.Coins,
	}

	for _, txn := range transactions {
		if txn.TransactionType == "received" {
			response.CoinHistory.Received = append(response.CoinHistory.Received, struct {
				FromUser string `json:"fromUser"`
				Amount   int    `json:"amount"`
			}{
				FromUser: txn.CounterpartUser,
				Amount:   txn.Amount,
			})
		} else if txn.TransactionType == "sent" {
			response.CoinHistory.Sent = append(response.CoinHistory.Sent, struct {
				ToUser string `json:"toUser"`
				Amount int    `json:"amount"`
			}{
				ToUser: txn.CounterpartUser,
				Amount: txn.Amount,
			})
		}
	}

	for _, item := range inventory {
		merch, err := h.merchService.GetMerchByID(r.Context(), item.ItemID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching merch: %v", err), http.StatusInternalServerError)
			return
		}
		if merch == nil {
			http.Error(w, fmt.Sprintf("Merch with that id not found, id: %d", item.ID), http.StatusBadRequest)
			return
		}

		response.Inventory = append(response.Inventory, struct {
			Type     string `json:"type"`
			Quantity int    `json:"quantity"`
		}{
			Type:     merch.ItemName,
			Quantity: item.Quantity,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}
