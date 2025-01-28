package repository

import "finance-api/internal/model"

type UserRepository interface {
    CreateUser(balance float64) (*model.User, error)
    GetUserByID(userID int) (*model.User, error)  // Добавляем метод для получения пользователя по ID
    UpdateUser(user *model.User) error            // Добавляем метод для обновления пользователя
	TransferBalance(fromUserID, toUserID int, amount float64) (*model.User, *model.User, error)
	TransferMoney(fromUserID, toUserID int, amount float64) error // Добавляем новый метод
}
