package service

import (
	"finance-api/internal/model"
	"fmt"
)

type UserRepository interface {
    CreateUser(balance float64) (*model.User, error)
    GetUserByID(userID int) (*model.User, error)  // Добавляем метод для получения пользователя по ID
    UpdateUser(user *model.User) error            // Добавляем метод для обновления пользователя
	TransferBalance(fromUserID, toUserID int, amount float64) (*model.User, *model.User, error)
	TransferMoney(fromUserID, toUserID int, amount float64) error // Добавляем новый метод
}

type UserService struct {
    Repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(balance float64) (*model.User, error) {
    return s.Repo.CreateUser(balance)
}

func (s *UserService) AddBalance(userID int, amount float64) (*model.User, error) {
    // Получаем пользователя
    user, err := s.Repo.GetUserByID(userID)
    if err != nil {
        return nil, fmt.Errorf("could not get user: %v", err)
    }

    // Обновляем баланс
    user.Balance += amount

    // Обновляем пользователя в базе данных
    if err := s.Repo.UpdateUser(user); err != nil {
        return nil, fmt.Errorf("could not update user: %v", err)
    }

    return user, nil
}

// GetUserByID возвращает пользователя по его ID
func (s *UserService) GetUserByID(userID int) (*model.User, error) {
    user, err := s.Repo.GetUserByID(userID)
    if err != nil {
        return nil, fmt.Errorf("could not get user: %v", err)
    }
    return user, nil
}
