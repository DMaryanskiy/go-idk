package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/DMaryanskiy/go-idk/internal/domain"
	"github.com/DMaryanskiy/go-idk/pkg/database"
)

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	query := `
	INSERT INTO users (email, name)
	VALUES ($1, $2)
	RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Name).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	user := &domain.User{}
	query := `
	SELECT id, email, name, created_at, updated_at
	FROM users
	WHERE id = $1;`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	user := &domain.User{}
	query := `
	SELECT id, email, name, created_at, updated_at
	FROM users
	WHERE email = $1;`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context, limit, offset int) (users []domain.User, total int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()

	countQuery := "SELECT COUNT(*) FROM users;"

	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting users: %w", err)
	}

	query := "SELECT id, email, name, created_at, updated_at FROM users ORDER BY id LIMIT $1 OFFSET $2;"
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting users: %w", err)
	}
	defer func() {
		errRows := rows.Close()
		if errRows != nil {
			err = errRows
		}
	}()

	users = []domain.User{}
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, id int, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	query := `
	UPDATE users
	SET email = $1, name = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	RETURNING updated_at;`

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Name, id).Scan(&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user not found")
	}
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	user.ID = id
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1;`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deleted rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
