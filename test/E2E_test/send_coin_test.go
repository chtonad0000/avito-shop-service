//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const recipientUsername = "testrecipient"
const recipientPassword = "recipientpass"
const transferAmount = 50

func TestE2E_TransferCoins(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	authResp, err := client.Post(baseURL+"/auth", "application/json", bytes.NewReader([]byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, authResp.StatusCode)

	var authData map[string]string
	err = json.NewDecoder(authResp.Body).Decode(&authData)
	assert.NoError(t, err)
	senderToken := authData["token"]
	assert.NotEmpty(t, senderToken)

	reqInfo, err := http.NewRequest("GET", baseURL+"/info", nil)
	assert.NoError(t, err)
	reqInfo.Header.Set("Authorization", "Bearer "+senderToken)

	infoResp, err := client.Do(reqInfo)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoResp.StatusCode)

	var senderInfo map[string]interface{}
	err = json.NewDecoder(infoResp.Body).Decode(&senderInfo)
	assert.NoError(t, err)
	coinsFloat, ok := senderInfo["coins"].(float64)
	if !ok {
		t.Fatalf("invalid type for coins, expected float64, got %T", senderInfo["coins"])
	}
	initialSenderCoins := int(coinsFloat)

	authRespRecipient, err := client.Post(baseURL+"/auth", "application/json", bytes.NewReader([]byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, recipientUsername, recipientPassword))))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, authRespRecipient.StatusCode)

	var authDataRecipient map[string]string
	err = json.NewDecoder(authRespRecipient.Body).Decode(&authDataRecipient)
	assert.NoError(t, err)
	recipientToken := authDataRecipient["token"]
	assert.NotEmpty(t, recipientToken)

	reqInfoRecipient, err := http.NewRequest("GET", baseURL+"/info", nil)
	assert.NoError(t, err)
	reqInfoRecipient.Header.Set("Authorization", "Bearer "+recipientToken)

	infoRespRecipient, err := client.Do(reqInfoRecipient)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoRespRecipient.StatusCode)

	var recipientInfo map[string]interface{}
	err = json.NewDecoder(infoRespRecipient.Body).Decode(&recipientInfo)
	assert.NoError(t, err)
	coinsFloat, ok = recipientInfo["coins"].(float64)
	if !ok {
		t.Fatalf("invalid type for coins, expected float64, got %T", recipientInfo["coins"])
	}
	initialRecipientCoins := int(coinsFloat)

	transferReq := map[string]interface{}{
		"toUser": recipientUsername,
		"amount": transferAmount,
	}
	transferReqBody, err := json.Marshal(transferReq)
	assert.NoError(t, err)
	reqTransfer, err := http.NewRequest("POST", baseURL+"/sendCoin", bytes.NewReader(transferReqBody))
	assert.NoError(t, err)
	reqTransfer.Header.Set("Authorization", "Bearer "+senderToken)

	transferResp, err := client.Do(reqTransfer)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, transferResp.StatusCode)

	infoRespAfterSender, err := client.Do(reqInfo)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoRespAfterSender.StatusCode)

	var senderInfoAfterTransfer map[string]interface{}
	err = json.NewDecoder(infoRespAfterSender.Body).Decode(&senderInfoAfterTransfer)
	assert.NoError(t, err)
	coinsFloat, ok = senderInfoAfterTransfer["coins"].(float64)
	if !ok {
		t.Fatalf("invalid type for coins, expected float64, got %T", senderInfoAfterTransfer["coins"])
	}
	finalSenderCoins := int(coinsFloat)
	infoRespAfterRecipient, err := client.Do(reqInfoRecipient)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoRespAfterRecipient.StatusCode)

	var recipientInfoAfterTransfer map[string]interface{}
	err = json.NewDecoder(infoRespAfterRecipient.Body).Decode(&recipientInfoAfterTransfer)
	assert.NoError(t, err)
	coinsFloat, ok = recipientInfoAfterTransfer["coins"].(float64)
	if !ok {
		t.Fatalf("invalid type for coins, expected float64, got %T", recipientInfoAfterTransfer["coins"])
	}
	finalRecipientCoins := int(coinsFloat)
	assert.Equal(t, initialSenderCoins-transferAmount, finalSenderCoins, "Sender balance must be previous balance minus amount")

	assert.Equal(t, initialRecipientCoins+transferAmount, finalRecipientCoins, "Receiver balance must be previous balance plus amount")
}
