package repository

import (
	"context"
	"finance-api/internal/model"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionRepositoryImpl struct {
	DB *pgxpool.Pool
}

// NewTransactionRepository создает новый экземпляр TransactionRepositoryImpl
func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{
		DB: db,
	}
}

// FetchLastTransactions возвращает последние транзакции пользователя
func (r *TransactionRepositoryImpl) FetchLastTransactions(userID int, limit int) ([]model.Transaction, error) {
    query := `
        SELECT id, sender_id, receiver_id, amount, created_at
        FROM transactions
        WHERE sender_id = $1 OR receiver_id = $1
        ORDER BY created_at DESC
        LIMIT $2;
    `

    rows, err := r.DB.Query(context.Background(), query, userID, limit)
    if err != nil {
        return nil, fmt.Errorf("could not fetch transactions: %v", err)
    }
    defer rows.Close()

    var transactions []model.Transaction
    for rows.Next() {
        var t model.Transaction
        err := rows.Scan(&t.ID, &t.SenderID, &t.ReceiverID, &t.Amount, &t.CreatedAt) // Сканируем created_at как time.Time
        if err != nil {
            return nil, fmt.Errorf("could not scan transaction: %v", err)
        }
        transactions = append(transactions, t)
    }

    if rows.Err() != nil {
        return nil, fmt.Errorf("error after scanning rows: %v", rows.Err())
    }

    return transactions, nil
}
