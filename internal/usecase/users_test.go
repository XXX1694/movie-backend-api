package usecase_test

import (
	"fmt"
	"testing"

	"golang/internal/mocks"
	"golang/internal/usecase"
	"golang/pkg/modules"

	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	mockRepo.On("GetUsers", 10, 0).Return([]modules.User{
		{ID: 1, Name: "John Doe", Email: "john@test.com", Age: 25},
	}, nil)

	uc := usecase.NewUserUsecase(mockRepo)
	users, err := uc.GetUsers(10, 0)

	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "John Doe", users[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	mockRepo.On("GetUserByID", 1).Return(&modules.User{
		ID: 1, Name: "John Doe", Email: "john@test.com", Age: 25,
	}, nil)

	uc := usecase.NewUserUsecase(mockRepo)
	user, err := uc.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "John Doe", user.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	mockRepo.On("GetUserByID", 999).Return(nil, fmt.Errorf("user with id 999 not found"))

	uc := usecase.NewUserUsecase(mockRepo)
	user, err := uc.GetUserByID(999)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	newUser := modules.User{Name: "Alice", Email: "alice@test.com", Age: 30}
	mockRepo.On("CreateUser", newUser).Return(1, nil)

	uc := usecase.NewUserUsecase(mockRepo)
	id, err := uc.CreateUser(newUser)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	mockRepo.On("DeleteUser", 999).Return(int64(0), fmt.Errorf("user with id 999 does not exist"))

	uc := usecase.NewUserUsecase(mockRepo)
	rows, err := uc.DeleteUser(999)

	assert.Error(t, err)
	assert.Equal(t, int64(0), rows)
	mockRepo.AssertExpectations(t)
}
