package movies

import (
	"database/sql"
	"errors"
	"fmt"

	_postgres "golang/internal/repository/_postgres"
	"golang/pkg/modules"
)

type Repository struct {
	db *_postgres.Dialect
}

func NewMovieRepository(db *_postgres.Dialect) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMovies(limit, offset int) ([]modules.Movie, error) {
	var movies []modules.Movie
	err := r.db.DB.Select(&movies,
		"SELECT id, title, description, year, rating, created_at FROM movies WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *Repository) GetMovieByID(id int) (*modules.Movie, error) {
	var movie modules.Movie
	err := r.db.DB.Get(&movie,
		"SELECT id, title, description, year, rating, created_at FROM movies WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("movie with id %d not found", id)
		}
		return nil, err
	}
	return &movie, nil
}

func (r *Repository) CreateMovie(movie modules.Movie) (int, error) {
	var id int
	err := r.db.DB.QueryRow(
		"INSERT INTO movies (title, description, year, rating) VALUES ($1, $2, $3, $4) RETURNING id",
		movie.Title, movie.Description, movie.Year, movie.Rating,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create movie: %w", err)
	}
	return id, nil
}

func (r *Repository) UpdateMovie(id int, movie modules.Movie) error {
	result, err := r.db.DB.Exec(
		"UPDATE movies SET title=$1, description=$2, year=$3, rating=$4 WHERE id=$5 AND deleted_at IS NULL",
		movie.Title, movie.Description, movie.Year, movie.Rating, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("movie with id %d does not exist", id)
	}
	return nil
}

func (r *Repository) DeleteMovie(id int) (int64, error) {
	result, err := r.db.DB.Exec(
		"UPDATE movies SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete movie: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return 0, fmt.Errorf("movie with id %d does not exist or already deleted", id)
	}
	return rows, nil
}
