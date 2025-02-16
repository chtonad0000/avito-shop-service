//go:build unit
// +build unit

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/avito-shop-service/internal/handlers"
	"github.com/avito-shop-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Auth_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserService.On("Authenticate", mock.Anything, "testuser", "password123").Return("test-token", nil)
	mockUserService.On("GetUserByUsername", mock.Anything, "testuser").Return(&models.User{Username: "testuser"}, nil)

	handler := handlers.NewUserHandler(mockUserService)
	reqBody := []byte(`{"username":"testuser","password":"password123"}`)
	req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(reqBody))
	rr := httptest.NewRecorder()

	handler.Auth(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "test-token", response["token"])
}

func TestUserHandler_Auth_Error_InvalidInput(t *testing.T) {
	mockUserService := new(MockUserService)

	handler := handlers.NewUserHandler(mockUserService)
	req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader([]byte(`{`)))
	rr := httptest.NewRecorder()

	handler.Auth(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandler_Auth_Error_CreatingUser(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserService.On("GetUserByUsername", mock.Anything, "testuser").Return((*models.User)(nil), nil)
	mockUserService.On("CreateUser", mock.Anything, "testuser", "password123").Return((*models.User)(nil), fmt.Errorf("error"))

	handler := handlers.NewUserHandler(mockUserService)
	reqBody := []byte(`{"username":"testuser","password":"password123"}`)
	req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(reqBody))
	rr := httptest.NewRecorder()
	handler.Auth(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUserHandler_Auth_Error_GetUserByUsername(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserService.On("GetUserByUsername", mock.Anything, "testuser").Return((*models.User)(nil), fmt.Errorf("database error"))

	handler := handlers.NewUserHandler(mockUserService)
	reqBody := []byte(`{"username":"testuser","password":"password123"}`)
	req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(reqBody))
	rr := httptest.NewRecorder()

	handler.Auth(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

}

func TestUserHandler_Auth_Error_Authenticate(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserService.On("GetUserByUsername", mock.Anything, "testuser").Return(&models.User{Username: "testuser"}, nil)
	mockUserService.On("Authenticate", mock.Anything, "testuser", "password123").Return("", fmt.Errorf("invalid password"))

	handler := handlers.NewUserHandler(mockUserService)
	reqBody := []byte(`{"username":"testuser","password":"password123"}`)
	req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(reqBody))
	rr := httptest.NewRecorder()

	handler.Auth(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
