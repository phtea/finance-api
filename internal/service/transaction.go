package service

import "finance-api/internal/model"

// TransactionService определяет интерфейс для работы с транзакциями
type TransactionService interface {
    // GetLastTransactions возвращает последние транзакции пользователя
    GetLastTransactions(userID int) ([]model.Transaction, error)
}
