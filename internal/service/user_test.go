package service

import (
	"errors"
	"finance-api/internal/model"
	"finance-api/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUserByID(t *testing.T) {
    mockRepo := new(mocks.MockUserRepository)
    userService := NewUserService(mockRepo)

    t.Run("Success", func(t *testing.T) {
        // Очистка моков перед настройкой
        mockRepo.ExpectedCalls = nil

        // Настройка мока
        mockRepo.On("GetUserByID", 1).Return(&model.User{ID: 1, Balance: 100.0}, nil)

        // Выполнение теста
        user, err := userService.GetUserByID(1)

        // Проверка результатов
        assert.NoError(t, err)
        assert.Equal(t, 1, user.ID)
        assert.Equal(t, 100.0, user.Balance)
        mockRepo.AssertExpectations(t)
    })

    t.Run("User Not Found", func(t *testing.T) {
        // Очистка моков перед настройкой
        mockRepo.ExpectedCalls = nil

        // Настройка мока
        mockRepo.On("GetUserByID", 1).Return(nil, errors.New("user not found"))

        // Выполнение теста
        user, err := userService.GetUserByID(1)

        // Проверка результатов
        assert.Error(t, err)
        assert.Nil(t, user)
        mockRepo.AssertExpectations(t)
    })
}

func TestAddBalance(t *testing.T) {
    // Создаем моковый репозиторий
    mockRepo := mocks.NewMockUserRepository()
    userService := &UserService{Repo: mockRepo}

    // Настройка поведения для GetUserByID и UpdateUser
    mockRepo.On("GetUserByID", 1).Return(&model.User{ID: 1, Balance: 100.00}, nil)
    mockRepo.On("UpdateUser", &model.User{ID: 1, Balance: 150.00}).Return(nil)

    // Тестируем AddBalance
    user, err := userService.AddBalance(1, 50.00)
    if err != nil {
        t.Fatalf("Error occurred: %v", err)
    }

    // Проверяем, что баланс обновился корректно
    assert.Equal(t, 150.00, user.Balance, "Баланс должен быть обновлен на 50.00")

    // Убеждаемся, что методы были вызваны
    mockRepo.AssertCalled(t, "GetUserByID", 1)
    mockRepo.AssertCalled(t, "UpdateUser", &model.User{ID: 1, Balance: 150.00})
}

