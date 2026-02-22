package users

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	_postgres "golang/internal/repository/_postgres"
	"golang/pkg/modules"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{db: db, executionTimeout: time.Second * 5}
}

func (r *Repository) GetUsers(limit, offset int) ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users,
		"SELECT id, name, email, age, created_at FROM users WHERE deleted_at IS NULL LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user,
		"SELECT id, name, email, age, created_at FROM users WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) DeleteUser(id int) (int64, error) {
	result, err := r.db.DB.Exec(
		"UPDATE users SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return 0, fmt.Errorf("user with id %d does not exist or already deleted", id)
	}
	return rows, nil
}

func (r *Repository) CreateUser(user modules.User) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	var id int
	err = r.db.DB.QueryRow(
		"INSERT INTO users (name, email, age, password) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, user.Age, string(hashedPassword),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (r *Repository) UpdateUser(id int, user modules.User) error {
	result, err := r.db.DB.Exec(
		"UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4 AND deleted_at IS NULL",
		user.Name, user.Email, user.Age, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user with id %d does not exist", id)
	}
	return nil
}

func (r *Repository) CreateUserWithAudit(user modules.User) (int, error) {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRow(
		"INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Age,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	_, err = tx.Exec(
		"INSERT INTO audit_log (user_id, action) VALUES ($1, $2)",
		id, "CREATE",
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create audit log: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}
