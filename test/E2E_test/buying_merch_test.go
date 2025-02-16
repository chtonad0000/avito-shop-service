//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080/api"
const itemName = "pen"
const merchPrice = 10
const username = "testuser"
const password = "pass"

func TestE2E_BuyMerch(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	authResp, err := client.Post(baseURL+"/auth", "application/json", bytes.NewReader([]byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))))
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

	reqBuy, err := http.NewRequest("GET", fmt.Sprintf("%s/buy/%s", baseURL, itemName), nil)
	assert.NoError(t, err)
	reqBuy.Header.Set("Authorization", "Bearer "+token)

	buyResp, err := client.Do(reqBuy)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, buyResp.StatusCode)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("error closing response body")
			return
		}
	}(buyResp.Body)
	infoResp2, err := client.Do(reqInfo)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoResp2.StatusCode)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("error closing response body")
			return
		}
	}(infoResp2.Body)
	var userInfo2 map[string]interface{}
	err = json.NewDecoder(infoResp2.Body).Decode(&userInfo2)
	assert.NoError(t, err)

	coinsFloat, ok = userInfo2["coins"].(float64)
	if !ok {
		t.Fatalf("invalid type for coins, expected float64, got %T", userInfo["coins"])
	}
	finalCoins := int(coinsFloat)
	assert.Equal(t, initialCoins-merchPrice, finalCoins, "Balance must be equal previous balance minus merch price")

	inventory, ok := userInfo2["inventory"].([]interface{})
	assert.True(t, ok, "Purchase must be array")

	itemFound := false
	for _, merch := range inventory {
		item, ok := merch.(map[string]interface{})
		if ok && item["type"] == itemName {
			itemFound = true
			break
		}
	}

	assert.True(t, itemFound, "Merch must be in inventory")
}
