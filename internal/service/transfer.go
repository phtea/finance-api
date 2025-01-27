package service

import (
	"fmt"
	"finance-api/internal/model"
)

func (s *UserService) TransferBalance(fromUserID, toUserID int, amount float64) (*model.User, *model.User, error) {
    // Проверка, что пользователь не переводит деньги сам себе
    if fromUserID == toUserID {
        return nil, nil, fmt.Errorf("cannot transfer to the same user")
    }

    // Получаем информацию об отправителе и получателе
    fromUser, err := s.Repo.GetUserByID(fromUserID)
    if err != nil {
        return nil, nil, fmt.Errorf("could not get from user: %v", err)
    }

    toUser, err := s.Repo.GetUserByID(toUserID)
    if err != nil {
        return nil, nil, fmt.Errorf("could not get to user: %v", err)
    }

    // Проверяем, что у отправителя достаточно средств
    if fromUser.Balance < amount {
        return nil, nil, fmt.Errorf("insufficient balance")
    }

    // Выполняем перевод
    if err := s.Repo.TransferMoney(fromUserID, toUserID, amount); err != nil {
        return nil, nil, fmt.Errorf("could not transfer money: %v", err)
    }

    // Обновляем балансы пользователей
    fromUser.Balance -= amount
    toUser.Balance += amount

    return fromUser, toUser, nil
}
