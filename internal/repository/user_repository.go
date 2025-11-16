package repository

import (
	"database/sql"
	"fmt"

	"github.com/DMaryanskiy/go-idk/internal/domain"
	"github.com/DMaryanskiy/go-idk/pkg/database"
)

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
	INSERT INTO users (name, email)
	($1, $2)
	RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, user.Name, user.Email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `
	SELECT id, name, email, created_at, updated_at
	FROM users
	WHERE id = $1;`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
	SELECT id, name, email, created_at, updated_at
	FROM users
	WHERE email = $1;`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetAll(limit, offset int) ([]domain.User, int, error) {
	var total int
	countQuery := "SELECT COUNT(*) FROM users;"

	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting users: %w", err)
	}

	query := "SELECT id, name, email, created_at, updated_at FROM users;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) Update(id int, user *domain.User) error {
	query := `
	UPDATE users
	SET email = $1, name = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	RETURNING updated_at;`

	err := r.db.QueryRow(query, user.Email, user.Name, id).Scan(&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1;`
	result, err := r.db.Exec(query, id)
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
