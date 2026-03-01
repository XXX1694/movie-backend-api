package usecase

import (
	"golang/internal/repository"
	"golang/pkg/modules"
)

type MovieUsecase struct {
	repo repository.MovieRepository
}

func NewMovieUsecase(repo repository.MovieRepository) *MovieUsecase {
	return &MovieUsecase{repo: repo}
}

func (u *MovieUsecase) GetMovies(limit, offset int) ([]modules.Movie, error) {
	return u.repo.GetMovies(limit, offset)
}

func (u *MovieUsecase) GetMovieByID(id int) (*modules.Movie, error) {
	return u.repo.GetMovieByID(id)
}

func (u *MovieUsecase) CreateMovie(movie modules.Movie) (int, error) {
	return u.repo.CreateMovie(movie)
}

func (u *MovieUsecase) UpdateMovie(id int, movie modules.Movie) error {
	return u.repo.UpdateMovie(id, movie)
}

func (u *MovieUsecase) DeleteMovie(id int) (int64, error) {
	return u.repo.DeleteMovie(id)
}

type MovieUsecaseInterface interface {
	GetMovies(limit, offset int) ([]modules.Movie, error)
	GetMovieByID(id int) (*modules.Movie, error)
	CreateMovie(movie modules.Movie) (int, error)
	UpdateMovie(id int, movie modules.Movie) error
	DeleteMovie(id int) (int64, error)
}
