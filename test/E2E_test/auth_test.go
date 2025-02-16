//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

const usernameAuth = "authuser"
const passwordAuth = "auth"

func TestAuthenticate(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	authResp, err := client.Post(baseURL+"/auth", "application/json", bytes.NewReader([]byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, usernameAuth, passwordAuth))))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, authResp.StatusCode)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("error closing response body")
			return
		}
	}(authResp.Body)

	var authData map[string]string
	err = json.NewDecoder(authResp.Body).Decode(&authData)
	assert.NoError(t, err)
	token := authData["token"]
	assert.NotEmpty(t, token)

	reqInfo, err := http.NewRequest("GET", baseURL+"/info", nil)
	assert.NoError(t, err)
	reqInfo.Header.Set("Authorization", "Bearer "+token)

	infoResp, err := client.Do(reqInfo)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoResp.StatusCode)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("error closing response body")
			return
		}
	}(infoResp.Body)
	var userInfo map[string]interface{}
	err = json.NewDecoder(infoResp.Body).Decode(&userInfo)
	assert.NoError(t, err)
	coinsFloat, ok := userInfo["coins"].(float64)
	if !ok {
		t.Fatalf("invalid type for coins, expected float64, got %T", userInfo["coins"])
	}
	initialCoins := int(coinsFloat)

	assert.Equal(t, 1000, initialCoins, "Must be 1000 coins after registration")
}
