package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "finance-api/internal/service"
)

type TransactionHandler struct {
    service service.TransactionService // Используем интерфейс, а не конкретную реализацию
}

// NewTransactionHandler принимает интерфейс TransactionService.
func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
    return &TransactionHandler{service: s}
}

// GetLastTransactions handles the request to get the last 10 transactions of a user.
func (h *TransactionHandler) GetLastTransactions(c *gin.Context) {
    // Извлекаем и проверяем user_id из query params
    userIDParam := c.Query("user_id")
    if userIDParam == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
        return
    }

    userID, err := strconv.Atoi(userIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
        return
    }

    // Получаем последние транзакции
    transactions, err := h.service.GetLastTransactions(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Возвращаем результат
    c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
