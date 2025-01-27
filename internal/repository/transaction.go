package repository

import "finance-api/internal/model"

// TransactionRepository определяет методы для работы с транзакциями
type TransactionRepository interface {
    // FetchLastTransactions возвращает последние транзакции пользователя
    FetchLastTransactions(userID int, limit int) ([]model.Transaction, error)
}
