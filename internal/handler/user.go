package handler

import (
	"finance-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    Service *service.UserService
}

// GetUserByID обрабатывает запрос на получение пользователя по ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
    // Извлекаем ID пользователя из параметров запроса
    userIDStr := c.Param("id")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    // Получаем пользователя по ID
    user, err := h.Service.GetUserByID(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Возвращаем данные пользователя
    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var userInput struct {
        Balance float64 `json:"balance"`
    }
    if err := c.BindJSON(&userInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    user, err := h.Service.CreateUser(userInput.Balance)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) AddBalance(c *gin.Context) {
    var input struct {
        UserID int     `json:"user_id"`
        Amount float64 `json:"amount"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    user, err := h.Service.AddBalance(input.UserID, input.Amount)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": user.ID, "balance": user.Balance})
}
