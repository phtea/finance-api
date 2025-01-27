package service

import (
    "errors"
    "finance-api/internal/model"
    "finance-api/internal/repository/mocks"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestTransactionService_GetLastTransactions(t *testing.T) {
    mockRepo := new(mocks.TransactionRepository)
    service := NewTransactionService(mockRepo)

    t.Run("Success", func(t *testing.T) {
        mockRepo.ExpectedCalls = nil // Сбрасываем настройки мока

        // Преобразуем строки в time.Time
        createdAt1, _ := time.Parse(time.RFC3339, "2025-01-26T10:00:00Z")
        createdAt2, _ := time.Parse(time.RFC3339, "2025-01-26T11:00:00Z")

        mockTransactions := []model.Transaction{
            {ID: 1, SenderID: 1, ReceiverID: 2, Amount: 100, CreatedAt: createdAt1},
            {ID: 2, SenderID: 1, ReceiverID: 3, Amount: 200, CreatedAt: createdAt2},
        }
        mockRepo.On("FetchLastTransactions", 1, 10).Return(mockTransactions, nil)

        transactions, err := service.GetLastTransactions(1)

        assert.NoError(t, err)
        assert.Equal(t, mockTransactions, transactions)
        mockRepo.AssertCalled(t, "FetchLastTransactions", 1, 10)
    })

    t.Run("Invalid User ID", func(t *testing.T) {
        transactions, err := service.GetLastTransactions(-1)

        assert.Error(t, err)
        assert.Nil(t, transactions)
    })

    t.Run("Repository Error", func(t *testing.T) {
        mockRepo.ExpectedCalls = nil // Сбрасываем настройки мока
        mockRepo.On("FetchLastTransactions", 1, 10).Return(nil, errors.New("db error"))

        transactions, err := service.GetLastTransactions(1)

        assert.Error(t, err)
        assert.Nil(t, transactions) // Ожидаем, что результат будет nil при ошибке
        mockRepo.AssertCalled(t, "FetchLastTransactions", 1, 10)
    })
}
