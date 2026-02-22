package usecase

import (
	"context"
	"golang/internal/cache"
	"golang/internal/repository"
	"golang/pkg/modules"
)

type UserUsecase struct {
	repo  repository.UserRepository
	cache *cache.RedisCache
}

func NewUserUsecase(repo repository.UserRepository, cache *cache.RedisCache) *UserUsecase {
	return &UserUsecase{repo: repo, cache: cache}
}

func (u *UserUsecase) GetUsers(limit, offset int) ([]modules.User, error) {
	return u.repo.GetUsers(limit, offset)
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	ctx := context.Background()

	// Проверяем кэш
	cached, err := u.cache.GetUser(ctx, id)
	if err == nil {
		return cached, nil
	}

	// Берём из БД
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	u.cache.SetUser(ctx, user)
	return user, nil
}

func (u *UserUsecase) CreateUser(user modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

func (u *UserUsecase) CreateUserWithAudit(user modules.User) (int, error) {
	return u.repo.CreateUserWithAudit(user)
}

func (u *UserUsecase) UpdateUser(id int, user modules.User) error {
	// Инвалидируем кэш
	u.cache.DeleteUser(context.Background(), id)
	return u.repo.UpdateUser(id, user)
}

func (u *UserUsecase) DeleteUser(id int) (int64, error) {
	// Инвалидируем кэш
	u.cache.DeleteUser(context.Background(), id)
	return u.repo.DeleteUser(id)
}

type UserUsecaseInterface interface {
	GetUsers(limit, offset int) ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user modules.User) (int, error)
	CreateUserWithAudit(user modules.User) (int, error)
	UpdateUser(id int, user modules.User) error
	DeleteUser(id int) (int64, error)
}
