package usecase

import (
	"golang/internal/repository"
	"golang/pkg/modules"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers(limit, offset int) ([]modules.User, error) {
	return u.repo.GetUsers(limit, offset)
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUsecase) CreateUser(user modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

func (u *UserUsecase) UpdateUser(id int, user modules.User) error {
	return u.repo.UpdateUser(id, user)
}

func (u *UserUsecase) DeleteUser(id int) (int64, error) {
	return u.repo.DeleteUser(id)
}

func (u *UserUsecase) CreateUserWithAudit(user modules.User) (int, error) {
	return u.repo.CreateUserWithAudit(user)
}

type UserUsecaseInterface interface {
	GetUsers(limit, offset int) ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user modules.User) (int, error)
	CreateUserWithAudit(user modules.User) (int, error)
	UpdateUser(id int, user modules.User) error
	DeleteUser(id int) (int64, error)
}
