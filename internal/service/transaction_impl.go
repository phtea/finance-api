package service

import (
    "errors"
    "finance-api/internal/model"
    "finance-api/internal/repository"
)

// TransactionServiceImpl реализует интерфейс TransactionService.
type TransactionServiceImpl struct {
    TransactionRepo repository.TransactionRepository
}

// NewTransactionService создает новый экземпляр TransactionServiceImpl.
func NewTransactionService(transactionRepo repository.TransactionRepository) *TransactionServiceImpl {
    return &TransactionServiceImpl{
        TransactionRepo: transactionRepo,
    }
}

// GetLastTransactions возвращает последние 10 транзакций пользователя.
func (s *TransactionServiceImpl) GetLastTransactions(userID int) ([]model.Transaction, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user ID")
    }

    transactions, err := s.TransactionRepo.FetchLastTransactions(userID, 10)
    if err != nil {
        return nil, err
    }

    if transactions == nil {
        return []model.Transaction{}, nil
    }

    return transactions, nil
}
