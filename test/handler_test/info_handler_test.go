//go:build unit
// +build unit

package handler

import (
	"errors"
	"github.com/avito-shop-service/internal/handlers"
	"github.com/avito-shop-service/internal/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInformationHandler_GetInfo_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockMerchService := new(MockMerchService)
	mockInventoryService := new(MockInventoryService)
	mockTransactionService := new(MockTransactionService)

	handler := handlers.NewInformationHandler(mockUserService, mockMerchService, mockInventoryService, mockTransactionService)

	req := httptest.NewRequest("GET", "/info", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	user := &models.User{ID: 1, Username: "testuser", Coins: 500}
	transactions := []models.CoinTransaction{
		{TransactionType: "received", CounterpartUser: "user1", Amount: 100},
		{TransactionType: "sent", CounterpartUser: "user2", Amount: 50},
	}
	inventory := []models.Inventory{
		{ItemID: 1, Quantity: 3},
		{ItemID: 2, Quantity: 1},
	}
	merch1 := &models.Merch{ID: 1, ItemName: "T-shirt"}
	merch2 := &models.Merch{ID: 2, ItemName: "Mug"}

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(user, nil)
	mockTransactionService.On("GetTransactionsByUserId", req.Context(), user.ID).Return(transactions, nil)
	mockInventoryService.On("GetInventoryByUserID", req.Context(), user.ID).Return(inventory, nil)
	mockMerchService.On("GetMerchByID", req.Context(), int64(1)).Return(merch1, nil)
	mockMerchService.On("GetMerchByID", req.Context(), int64(2)).Return(merch2, nil)

	handler.GetInfo(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"coins": 500,
		"inventory": [
			{"type": "T-shirt", "quantity": 3},
			{"type": "Mug", "quantity": 1}
		],
		"coinHistory": {
			"received": [{"fromUser": "user1", "amount": 100}],
			"sent": [{"toUser": "user2", "amount": 50}]
		}
	}`, w.Body.String())

	mockUserService.AssertExpectations(t)
	mockTransactionService.AssertExpectations(t)
	mockInventoryService.AssertExpectations(t)
	mockMerchService.AssertExpectations(t)
}

func TestInformationHandler_GetInfo_Unauthorized(t *testing.T) {
	handler := handlers.NewInformationHandler(nil, nil, nil, nil)

	req := httptest.NewRequest("GET", "/info", nil)
	w := httptest.NewRecorder()

	handler.GetInfo(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not authorized")
}

func TestInformationHandler_GetInfo_UserFetchError(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := handlers.NewInformationHandler(mockUserService, nil, nil, nil)

	req := httptest.NewRequest("GET", "/info", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return((*models.User)(nil), errors.New("DB error"))

	handler.GetInfo(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error fetching user")
}

func TestInformationHandler_GetInfo_TransactionFetchError(t *testing.T) {
	mockUserService := new(MockUserService)
	mockTransactionService := new(MockTransactionService)
	handler := handlers.NewInformationHandler(mockUserService, nil, nil, mockTransactionService)

	req := httptest.NewRequest("GET", "/info", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	user := &models.User{ID: 1, Username: "testuser", Coins: 500}

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(user, nil)
	mockTransactionService.On("GetTransactionsByUserId", req.Context(), user.ID).Return(([]models.CoinTransaction)(nil), errors.New("DB error"))

	handler.GetInfo(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error fetching transactions")
}

func TestInformationHandler_GetInfo_InventoryFetchError(t *testing.T) {
	mockUserService := new(MockUserService)
	mockTransactionService := new(MockTransactionService)
	mockInventoryService := new(MockInventoryService)
	handler := handlers.NewInformationHandler(mockUserService, nil, mockInventoryService, mockTransactionService)

	req := httptest.NewRequest("GET", "/info", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	user := &models.User{ID: 1, Username: "testuser", Coins: 500}

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(user, nil)
	mockTransactionService.On("GetTransactionsByUserId", req.Context(), user.ID).Return([]models.CoinTransaction{}, nil)
	mockInventoryService.On("GetInventoryByUserID", req.Context(), user.ID).Return(([]models.Inventory)(nil), errors.New("DB error"))

	handler.GetInfo(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error fetching inventory")
}

func TestInformationHandler_GetInfo_MerchFetchError(t *testing.T) {
	mockUserService := new(MockUserService)
	mockTransactionService := new(MockTransactionService)
	mockInventoryService := new(MockInventoryService)
	mockMerchService := new(MockMerchService)
	handler := handlers.NewInformationHandler(mockUserService, mockMerchService, mockInventoryService, mockTransactionService)

	req := httptest.NewRequest("GET", "/info", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	user := &models.User{ID: 1, Username: "testuser", Coins: 500}
	inventory := []models.Inventory{
		{ItemID: 1, Quantity: 3},
	}

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(user, nil)
	mockTransactionService.On("GetTransactionsByUserId", req.Context(), user.ID).Return([]models.CoinTransaction{}, nil)
	mockInventoryService.On("GetInventoryByUserID", req.Context(), user.ID).Return(inventory, nil)
	mockMerchService.On("GetMerchByID", req.Context(), int64(1)).Return((*models.Merch)(nil), errors.New("DB error"))

	handler.GetInfo(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error fetching merch")
}

func TestInformationHandler_GetInfo_UserNotFound(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := handlers.NewInformationHandler(mockUserService, nil, nil, nil)

	req := httptest.NewRequest("GET", "/info", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return((*models.User)(nil), nil)

	handler.GetInfo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error finding user")
}

func TestInformationHandler_GetInfo_MerchNotFound(t *testing.T) {
	mockUserService := new(MockUserService)
	mockTransactionService := new(MockTransactionService)
	mockInventoryService := new(MockInventoryService)
	mockMerchService := new(MockMerchService)
	handler := handlers.NewInformationHandler(mockUserService, mockMerchService, mockInventoryService, mockTransactionService)

	req := httptest.NewRequest("GET", "/info", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	user := &models.User{ID: 1, Username: "testuser", Coins: 500}
	inventory := []models.Inventory{
		{ItemID: 1, Quantity: 3},
	}

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(user, nil)
	mockTransactionService.On("GetTransactionsByUserId", req.Context(), user.ID).Return([]models.CoinTransaction{}, nil)
	mockInventoryService.On("GetInventoryByUserID", req.Context(), user.ID).Return(inventory, nil)
	mockMerchService.On("GetMerchByID", req.Context(), int64(1)).Return((*models.Merch)(nil), nil)

	handler.GetInfo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Merch with that id not found")
}
