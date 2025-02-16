//go:build unit
// +build unit

package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avito-shop-service/internal/handlers"
	"github.com/avito-shop-service/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestBuyHandler_Buy_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockMerchService := new(MockMerchService)
	mockInventoryService := new(MockInventoryService)
	mockTransactionService := new(MockTransactionService)

	handler := handlers.NewBuyHandler(mockUserService, mockMerchService, mockInventoryService, mockTransactionService)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("item", "item1")
	req := httptest.NewRequest("POST", "/buy/item1", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	mockMerchService.On("GetMerchByName", req.Context(), "item1").Return(&models.Merch{ID: 1, ItemName: "item1", Price: 100}, nil)
	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(&models.User{ID: 1, Username: "testuser", Coins: 200}, nil)
	mockInventoryService.On("BuyItemToInventory", req.Context(), int64(1), int64(1), 1, 100).Return(nil)

	handler.Buy(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "purchase successful", response["message"])

	mockMerchService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
	mockInventoryService.AssertExpectations(t)
}

func TestBuyHandler_Buy_ItemNotFound(t *testing.T) {
	mockUserService := new(MockUserService)
	mockMerchService := new(MockMerchService)
	mockInventoryService := new(MockInventoryService)
	mockTransactionService := new(MockTransactionService)

	handler := handlers.NewBuyHandler(mockUserService, mockMerchService, mockInventoryService, mockTransactionService)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("item", "item1")
	req := httptest.NewRequest("POST", "/buy/item1", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	mockMerchService.On("GetMerchByName", req.Context(), "item1").Return((*models.Merch)(nil), nil)

	handler.Buy(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyHandler_Buy_NotEnoughCoins(t *testing.T) {
	mockUserService := new(MockUserService)
	mockMerchService := new(MockMerchService)
	mockInventoryService := new(MockInventoryService)
	mockTransactionService := new(MockTransactionService)

	handler := handlers.NewBuyHandler(mockUserService, mockMerchService, mockInventoryService, mockTransactionService)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("item", "item1")
	req := httptest.NewRequest("POST", "/buy/item1", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	mockMerchService.On("GetMerchByName", req.Context(), "item1").Return(&models.Merch{ID: 1, ItemName: "item1", Price: 100}, nil)
	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(&models.User{ID: 1, Username: "testuser", Coins: 50}, nil)

	handler.Buy(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyHandler_Buy_ErrorFetchingUser(t *testing.T) {
	mockUserService := new(MockUserService)
	mockMerchService := new(MockMerchService)
	mockInventoryService := new(MockInventoryService)
	mockTransactionService := new(MockTransactionService)

	handler := handlers.NewBuyHandler(mockUserService, mockMerchService, mockInventoryService, mockTransactionService)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("item", "item1")
	req := httptest.NewRequest("POST", "/buy/item1", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	mockMerchService.On("GetMerchByName", req.Context(), "item1").Return(&models.Merch{ID: 1, ItemName: "item1", Price: 100}, nil)
	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return((*models.User)(nil), errors.New("error fetching user"))

	handler.Buy(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBuyHandler_Buy_InventoryUpdateError(t *testing.T) {
	mockUserService := new(MockUserService)
	mockMerchService := new(MockMerchService)
	mockInventoryService := new(MockInventoryService)
	mockTransactionService := new(MockTransactionService)

	handler := handlers.NewBuyHandler(mockUserService, mockMerchService, mockInventoryService, mockTransactionService)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("item", "item1")

	req := httptest.NewRequest("POST", "/buy/item1", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	mockMerchService.On("GetMerchByName", req.Context(), "item1").Return(&models.Merch{ID: 1, ItemName: "item1", Price: 100}, nil)
	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(&models.User{ID: 1, Username: "testuser", Coins: 200}, nil)
	mockInventoryService.On("BuyItemToInventory", req.Context(), int64(1), int64(1), 1, 100).Return(errors.New("DB error"))

	handler.Buy(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error updating inventory")

	mockMerchService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
	mockInventoryService.AssertExpectations(t)
}

func TestBuyHandler_Buy_UserNotAuthorized(t *testing.T) {
	handler := handlers.NewBuyHandler(nil, nil, nil, nil)

	req := httptest.NewRequest("POST", "/buy/item1", nil)
	w := httptest.NewRecorder()

	handler.Buy(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not authorized")
}

func TestBuyHandler_Buy_ItemNameMissing(t *testing.T) {
	handler := handlers.NewBuyHandler(nil, nil, nil, nil)

	req := httptest.NewRequest("POST", "/buy/", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	handler.Buy(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "item name is required")
}

func TestBuyHandler_Buy_MerchFetchError(t *testing.T) {
	mockUserService := new(MockUserService)
	mockMerchService := new(MockMerchService)
	mockInventoryService := new(MockInventoryService)
	mockTransactionService := new(MockTransactionService)

	handler := handlers.NewBuyHandler(mockUserService, mockMerchService, mockInventoryService, mockTransactionService)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("item", "item1")

	req := httptest.NewRequest("POST", "/buy/item1", nil)
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	mockMerchService.On("GetMerchByName", req.Context(), "item1").Return((*models.Merch)(nil), errors.New("DB error"))

	handler.Buy(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error fetching merch")

	mockMerchService.AssertExpectations(t)
}
