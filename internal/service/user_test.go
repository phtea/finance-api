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

func TestUserService_AddBalance(t *testing.T) {
	// Создаем моковый репозиторий
	mockRepo := mocks.NewMockUserRepository()
	userService := NewUserService(mockRepo)
	// userService := &UserService{Repo: mockRepo}

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

func TestUserService_TransferBalance(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	service := NewUserService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil

		// Mock users
		fromUser := &model.User{ID: 1, Balance: 200}
		toUser := &model.User{ID: 2, Balance: 50}

		// Mock repository calls
		mockRepo.On("GetUserByID", 1).Return(fromUser, nil)
		mockRepo.On("GetUserByID", 2).Return(toUser, nil)
		mockRepo.On("TransferMoney", 1, 2, 100.0).Return(nil)

		updatedFromUser, updatedToUser, err := service.TransferBalance(1, 2, 100.0)

		assert.NoError(t, err)
		assert.Equal(t, 100.0, updatedFromUser.Balance)
		assert.Equal(t, 150.0, updatedToUser.Balance)
		mockRepo.AssertCalled(t, "GetUserByID", 1)
		mockRepo.AssertCalled(t, "GetUserByID", 2)
		mockRepo.AssertCalled(t, "TransferMoney", 1, 2, 100.0)
	})

	t.Run("Same User Transfer", func(t *testing.T) {
		_, _, err := service.TransferBalance(1, 1, 100.0)
		assert.Error(t, err)
		assert.Equal(t, "cannot transfer to the same user", err.Error())
	})

	t.Run("Insufficient Balance", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil

		fromUser := &model.User{ID: 1, Balance: 50}
		mockRepo.On("GetUserByID", 1).Return(fromUser, nil)

		toUser := &model.User{ID: 2, Balance: 50}
		mockRepo.On("GetUserByID", 2).Return(toUser, nil)

		_, _, err := service.TransferBalance(1, 2, 100.0)

		assert.Error(t, err)
		assert.Equal(t, "insufficient balance", err.Error())
		mockRepo.AssertCalled(t, "GetUserByID", 1)
		mockRepo.AssertCalled(t, "GetUserByID", 2)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil

		fromUser := &model.User{ID: 1, Balance: 200}
		mockRepo.On("GetUserByID", 1).Return(fromUser, nil)

		toUser := &model.User{ID: 2, Balance: 50}
		mockRepo.On("GetUserByID", 2).Return(toUser, nil)

		mockRepo.On("TransferMoney", 1, 2, 100.0).Return(errors.New("db error"))

		_, _, err := service.TransferBalance(1, 2, 100.0)

		assert.Error(t, err)
		assert.Equal(t, "could not transfer money: db error", err.Error())
		mockRepo.AssertCalled(t, "GetUserByID", 1)
		mockRepo.AssertCalled(t, "GetUserByID", 2)
		mockRepo.AssertCalled(t, "TransferMoney", 1, 2, 100.0)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	service := NewUserService(mockRepo) // Передаём мок в сервис

	t.Run("Success", func(t *testing.T) {
		initialBalance := 100.0
		expectedUser := &model.User{ID: 1, Balance: initialBalance}

		// Настраиваем мок: возвращаем объект пользователя и nil в качестве ошибки
		mockRepo.On("CreateUser", initialBalance).Return(expectedUser, nil)

		// Вызываем метод CreateUser
		createdUser, err := service.CreateUser(initialBalance)

		// Проверяем, что ошибок нет и пользователь создан корректно
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, expectedUser, createdUser)

		// Проверяем, что метод мока был вызван с ожидаемым параметром
		mockRepo.AssertCalled(t, "CreateUser", initialBalance)
	})

	// Проверяем, что все ожидания мока выполнены
	mockRepo.AssertExpectations(t)
}
