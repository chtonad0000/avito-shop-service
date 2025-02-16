//go:build unit
// +build unit

package handler

import (
	"errors"
	"github.com/avito-shop-service/internal/handlers"
	"github.com/avito-shop-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTransactionHandler_SendCoin_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockTransactionService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockUserService, mockTransactionService)

	reqBody := `{"toUser": "recipient", "amount": 100}`
	req := httptest.NewRequest("POST", "/send-coin", strings.NewReader(reqBody))
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	user := &models.User{ID: 1, Username: "testuser", Coins: 500}

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(user, nil)
	mockTransactionService.On("CreateTransaction", req.Context(), mock.AnythingOfType("*models.CoinTransaction")).Return(nil)

	handler.SendCoin(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "send successful")
}

func TestTransactionHandler_SendCoin_Unauthorized(t *testing.T) {
	handler := handlers.NewTransactionHandler(nil, nil)

	req := httptest.NewRequest("POST", "/send-coin", nil)
	w := httptest.NewRecorder()

	handler.SendCoin(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not authorized")
}

func TestTransactionHandler_SendCoin_InvalidRequestBody(t *testing.T) {
	handler := handlers.NewTransactionHandler(nil, nil)

	req := httptest.NewRequest("POST", "/send-coin", strings.NewReader("{invalid_json"))
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	handler.SendCoin(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestTransactionHandler_SendCoin_InvalidToUserOrAmount(t *testing.T) {
	handler := handlers.NewTransactionHandler(nil, nil)

	reqBody := `{"toUser": "", "amount": 100}`
	req := httptest.NewRequest("POST", "/send-coin", strings.NewReader(reqBody))
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	handler.SendCoin(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid toUser or amount")
}

func TestTransactionHandler_SendCoin_UserNotFound(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := handlers.NewTransactionHandler(mockUserService, nil)

	reqBody := `{"toUser": "recipient", "amount": 100}`
	req := httptest.NewRequest("POST", "/send-coin", strings.NewReader(reqBody))
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return((*models.User)(nil), errors.New("user not found"))

	handler.SendCoin(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error fetching user")
}

func TestTransactionHandler_SendCoin_NotEnoughCoins(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := handlers.NewTransactionHandler(mockUserService, nil)

	reqBody := `{"toUser": "recipient", "amount": 1000}`
	req := httptest.NewRequest("POST", "/send-coin", strings.NewReader(reqBody))
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	user := &models.User{ID: 1, Username: "testuser", Coins: 500}

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(user, nil)

	handler.SendCoin(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "not enough money")
}

func TestTransactionHandler_SendCoin_CreateTransactionError(t *testing.T) {
	mockUserService := new(MockUserService)
	mockTransactionService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockUserService, mockTransactionService)

	reqBody := `{"toUser": "recipient", "amount": 100}`
	req := httptest.NewRequest("POST", "/send-coin", strings.NewReader(reqBody))
	req = req.WithContext(setEmployeeUsername(req.Context(), "testuser"))
	w := httptest.NewRecorder()

	user := &models.User{ID: 1, Username: "testuser", Coins: 500}

	mockUserService.On("GetUserByUsername", req.Context(), "testuser").Return(user, nil)
	mockTransactionService.On("CreateTransaction", req.Context(), mock.AnythingOfType("*models.CoinTransaction")).Return(errors.New("db error"))

	handler.SendCoin(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to send coins")
}
