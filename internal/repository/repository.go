package repository

import (
	_postgres "golang/internal/repository/_postgres"
	"golang/internal/repository/_postgres/movies"
	"golang/internal/repository/_postgres/users"
	"golang/pkg/modules"
)

type UserRepository interface {
	GetUsers(limit, offset int) ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user modules.User) (int, error)
	CreateUserWithAudit(user modules.User) (int, error)
	UpdateUser(id int, user modules.User) error
	DeleteUser(id int) (int64, error)
}

type MovieRepository interface {
	GetMovies(limit, offset int) ([]modules.Movie, error)
	GetMovieByID(id int) (*modules.Movie, error)
	CreateMovie(movie modules.Movie) (int, error)
	UpdateMovie(id int, movie modules.Movie) error
	DeleteMovie(id int) (int64, error)
}

type Repositories struct {
	UserRepository
	MovieRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository:  users.NewUserRepository(db),
		MovieRepository: movies.NewMovieRepository(db),
	}
}
