package mocks

import (
	"finance-api/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository - Mock implementation of UserRepository
// using testify's mock.Mock

type MockUserRepository struct {
    mock.Mock
}

// NewMockUserRepository creates a new instance of MockUserRepository
func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{}
}

// CreateUser - Mock implementation for creating a new user with a given balance
func (m *MockUserRepository) CreateUser(balance float64) (*model.User, error) {
    args := m.Called(balance)
    if args.Get(0) != nil {
        return args.Get(0).(*model.User), args.Error(1)
    }
    return nil, args.Error(1)
}

// GetUserByID - Mock implementation for retrieving a user by their ID
func (m *MockUserRepository) GetUserByID(userID int) (*model.User, error) {
    args := m.Called(userID)
    if args.Get(0) != nil {
        return args.Get(0).(*model.User), args.Error(1)
    }
    return nil, args.Error(1)
}

// UpdateUser - Mock implementation for updating an existing user
func (m *MockUserRepository) UpdateUser(user *model.User) error {
    args := m.Called(user)
    return args.Error(0)
}

// TransferBalance - Mock implementation for transferring balance between users
func (m *MockUserRepository) TransferBalance(fromUserID, toUserID int, amount float64) (*model.User, *model.User, error) {
    args := m.Called(fromUserID, toUserID, amount)
    if args.Get(0) != nil && args.Get(1) != nil {
        return args.Get(0).(*model.User), args.Get(1).(*model.User), args.Error(2)
    }
    return nil, nil, args.Error(2)
}

// TransferMoney - Mock implementation for transferring money and recording the transaction
func (m *MockUserRepository) TransferMoney(fromUserID, toUserID int, amount float64) error {
    args := m.Called(fromUserID, toUserID, amount)
    return args.Error(0)
}
