package handler

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
    FromUserID int     `json:"from_user_id"`
    ToUserID   int     `json:"to_user_id"`
    Amount     float64 `json:"amount"`
}

func (h *UserHandler) TransferBalance(c *gin.Context) {
    var input TransferRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    fromUser, toUser, err := h.Service.TransferBalance(input.FromUserID, input.ToUserID, input.Amount)
	fmt.Println("From user:", fromUser)
	fmt.Println("To user:", toUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "from_user": gin.H{"id": fromUser.ID, "balance": fromUser.Balance},
        "to_user":   gin.H{"id": toUser.ID, "balance": toUser.Balance},
    })
}
