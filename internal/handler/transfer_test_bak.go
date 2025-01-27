package handler

import (
	"bytes"
	"encoding/json"
	"finance-api/internal/repository/mocks"
	"finance-api/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTransferBalance(t *testing.T) {
    // Mock repository for user data
	mockRepo := mocks.NewMockUserRepository() // Это инициализирует репозиторий с пользователями
    userService := &service.UserService{Repo: mockRepo}

    // Create handler with the service
    handler := &UserHandler{Service: userService}

    // Input for the transfer request
    transferInput := struct {
        FromUserID int     `json:"from_user_id"`
        ToUserID   int     `json:"to_user_id"`
        Amount     float64 `json:"amount"`
    }{
        FromUserID: 1,    // Sender ID
        ToUserID:   2,    // Recipient ID
        Amount:     50.00, // Amount to transfer
    }

    // Marshal the input into JSON
    body, err := json.Marshal(transferInput)
    if err != nil {
        t.Fatalf("Could not marshal input data: %v", err)
    }

    // Create a request for the transfer endpoint
    req, err := http.NewRequest("POST", "/users/transfer", bytes.NewBuffer(body))
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Record the response
    w := httptest.NewRecorder()
    router := gin.Default()
    router.POST("/users/transfer", handler.TransferBalance)
    router.ServeHTTP(w, req)

    // Assert that the response code is OK (200)
    assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response contains the correct balance for both users
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	// Make sure the expected keys are present in the response
	fromUser, fromOk := response["from_user"].(map[string]interface{})
	toUser, toOk := response["to_user"].(map[string]interface{})

	if !fromOk || !toOk {
		t.Fatalf("Expected 'from_user' and 'to_user' in the response body")
	}

	// Assert the balances
	assert.NotNil(t, fromUser["id"])
	assert.NotNil(t, toUser["id"])
	assert.Equal(t, 100.00, fromUser["balance"])
	assert.Equal(t, 50.00, toUser["balance"])
}
