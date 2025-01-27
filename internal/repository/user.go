package repository

import (
    "context"
    "finance-api/internal/model"
    "fmt"
    "github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
    DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(balance float64) (*model.User, error) {
    var user model.User
    query := `INSERT INTO users (balance) VALUES ($1) RETURNING id, balance`
    err := r.DB.QueryRow(context.Background(), query, balance).Scan(&user.ID, &user.Balance)
    if err != nil {
        return nil, fmt.Errorf("could not create user: %v", err)
    }
    return &user, nil
}

func (r *UserRepository) GetUserByID(userID int) (*model.User, error) {
    var user model.User
    query := "SELECT id, balance FROM users WHERE id = $1"
    err := r.DB.QueryRow(context.Background(), query, userID).Scan(&user.ID, &user.Balance)
    if err != nil {
        return nil, fmt.Errorf("could not get user: %v", err)
    }
    return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
    // Начинаем транзакцию
    tx, err := r.DB.Begin(context.Background())
    if err != nil {
        return fmt.Errorf("could not begin transaction: %v", err)
    }
    defer tx.Rollback(context.Background()) // Откат транзакции в случае ошибки

    // Обновляем баланс пользователя
    query := "UPDATE users SET balance = $1 WHERE id = $2"
    _, err = tx.Exec(context.Background(), query, user.Balance, user.ID)
    if err != nil {
        return fmt.Errorf("could not update user balance: %v", err)
    }

    // Записываем транзакцию в таблицу transactions
    _, err = tx.Exec(context.Background(), `
        INSERT INTO transactions (sender_id, receiver_id, amount)
        VALUES ($1, $1, $2)
    `, user.ID, user.Balance)
    if err != nil {
        return fmt.Errorf("could not record transaction: %v", err)
    }

    // Фиксируем транзакцию
    if err := tx.Commit(context.Background()); err != nil {
        return fmt.Errorf("could not commit transaction: %v", err)
    }

    return nil
}

// TransferMoney выполняет перевод денег между пользователями и записывает транзакцию
func (r *UserRepository) TransferMoney(fromUserID, toUserID int, amount float64) error {
    // Начинаем транзакцию
    tx, err := r.DB.Begin(context.Background())
    if err != nil {
        return fmt.Errorf("could not begin transaction: %v", err)
    }
    defer tx.Rollback(context.Background()) // Откат транзакции в случае ошибки

    // Проверяем, что отправитель и получатель существуют
    var fromUserExists, toUserExists bool
    err = tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", fromUserID).Scan(&fromUserExists)
    if err != nil {
        return fmt.Errorf("could not check from user existence: %v", err)
    }
    err = tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", toUserID).Scan(&toUserExists)
    if err != nil {
        return fmt.Errorf("could not check to user existence: %v", err)
    }
    if !fromUserExists || !toUserExists {
        return fmt.Errorf("sender or receiver not found")
    }

    // Проверяем, что у отправителя достаточно средств
    var fromUserBalance float64
    err = tx.QueryRow(context.Background(), "SELECT balance FROM users WHERE id = $1", fromUserID).Scan(&fromUserBalance)
    if err != nil {
        return fmt.Errorf("could not get sender balance: %v", err)
    }
    if fromUserBalance < amount {
        return fmt.Errorf("insufficient funds")
    }

    // Обновляем баланс отправителя
    _, err = tx.Exec(context.Background(), "UPDATE users SET balance = balance - $1 WHERE id = $2", amount, fromUserID)
    if err != nil {
        return fmt.Errorf("could not deduct from sender: %v", err)
    }

    // Обновляем баланс получателя
    _, err = tx.Exec(context.Background(), "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toUserID)
    if err != nil {
        return fmt.Errorf("could not add to receiver: %v", err)
    }

    // Записываем транзакцию
    _, err = tx.Exec(context.Background(), `
        INSERT INTO transactions (sender_id, receiver_id, amount)
        VALUES ($1, $2, $3)
    `, fromUserID, toUserID, amount)
    if err != nil {
        return fmt.Errorf("could not record transaction: %v", err)
    }

    // Фиксируем транзакцию
    if err := tx.Commit(context.Background()); err != nil {
        return fmt.Errorf("could not commit transaction: %v", err)
    }

    return nil
}

func (r *UserRepository) TransferBalance(fromUserID, toUserID int, amount float64) (*model.User, *model.User, error) {
    tx, err := r.DB.Begin(context.Background())
    if err != nil {
        return nil, nil, fmt.Errorf("could not begin transaction: %v", err)
    }
    defer tx.Rollback(context.Background()) // Rollback in case of failure

    // Check if both users exist
    var fromUserExists, toUserExists bool
    err = tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", fromUserID).Scan(&fromUserExists)
    if err != nil {
        return nil, nil, fmt.Errorf("could not check from user existence: %v", err)
    }
    err = tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", toUserID).Scan(&toUserExists)
    if err != nil {
        return nil, nil, fmt.Errorf("could not check to user existence: %v", err)
    }
    if !fromUserExists || !toUserExists {
        return nil, nil, fmt.Errorf("sender or receiver not found")
    }

    // Check sender's balance
    var fromUserBalance float64
    err = tx.QueryRow(context.Background(), "SELECT balance FROM users WHERE id = $1", fromUserID).Scan(&fromUserBalance)
    if err != nil {
        return nil, nil, fmt.Errorf("could not get sender balance: %v", err)
    }
    if fromUserBalance < amount {
        return nil, nil, fmt.Errorf("insufficient funds")
    }

    // Deduct amount from sender
    _, err = tx.Exec(context.Background(), "UPDATE users SET balance = balance - $1 WHERE id = $2", amount, fromUserID)
    if err != nil {
        return nil, nil, fmt.Errorf("could not deduct from sender: %v", err)
    }

    // Add amount to receiver
    _, err = tx.Exec(context.Background(), "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toUserID)
    if err != nil {
        return nil, nil, fmt.Errorf("could not add to receiver: %v", err)
    }

    // Insert transaction records
    _, err = tx.Exec(context.Background(), `
        INSERT INTO transactions (sender_id, receiver_id, amount)
        VALUES ($1, $2, $3)
    `, fromUserID, toUserID, amount)
    if err != nil {
        return nil, nil, fmt.Errorf("could not record transaction: %v", err)
    }

    // Commit transaction
    if err := tx.Commit(context.Background()); err != nil {
        return nil, nil, fmt.Errorf("could not commit transaction: %v", err)
    }

    // Retrieve the updated users
    fromUser, err := r.GetUserByID(fromUserID)
    if err != nil {
        return nil, nil, fmt.Errorf("could not retrieve from user: %v", err)
    }

    toUser, err := r.GetUserByID(toUserID)
    if err != nil {
        return nil, nil, fmt.Errorf("could not retrieve to user: %v", err)
    }

    return fromUser, toUser, nil
}
