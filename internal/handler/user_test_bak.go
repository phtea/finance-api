package handler

import (
    "bytes"
    "encoding/json"
    "finance-api/internal/repository/mocks"
    "finance-api/internal/service"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateUser(t *testing.T) {
    // Инициализация мок-репозитория с 3 пользователями
    mockRepo := mocks.NewMockUserRepository()
    userService := &service.UserService{Repo: mockRepo}

    handler := &UserHandler{Service: userService}

    userInput := struct {
        Balance float64 `json:"balance"`
    }{
        Balance: 100.00,
    }

    body, err := json.Marshal(userInput)
    if err != nil {
        t.Fatalf("Could not marshal input data: %v", err)
    }

    req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router := gin.Default()
    router.POST("/users", handler.CreateUser)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var userResponse map[string]interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &userResponse); err != nil {
        t.Fatalf("Could not parse response: %v", err)
    }
    assert.NotNil(t, userResponse["id"])
    assert.Equal(t, userInput.Balance, userResponse["balance"])
}

func TestAddBalance(t *testing.T) {
    // Инициализация мок-репозитория
    mockRepo := mocks.NewMockUserRepository()
    userService := &service.UserService{Repo: mockRepo}

    // Создадим пользователя с ID 1 и проверим ошибку
    user, err := mockRepo.CreateUser(100.00)
    if err != nil {
        t.Fatalf("Could not create user: %v", err)
    }

    handler := &UserHandler{Service: userService}

    userInput := struct {
        UserID int     `json:"user_id"`
        Amount float64 `json:"amount"`
    }{
        UserID: user.ID,
        Amount: 50.00,
    }

    body, err := json.Marshal(userInput)
    if err != nil {
        t.Fatalf("Could not marshal input data: %v", err)
    }

    req, err := http.NewRequest("POST", "/users/balance", bytes.NewBuffer(body))
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router := gin.Default()
    router.POST("/users/balance", handler.AddBalance)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var userResponse map[string]interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &userResponse); err != nil {
        t.Fatalf("Could not parse response: %v", err)
    }
    assert.NotNil(t, userResponse["id"])
    assert.Equal(t, 150.00, userResponse["balance"])
}
