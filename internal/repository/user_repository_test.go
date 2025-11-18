package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DMaryanskiy/go-idk/internal/domain"
	"github.com/DMaryanskiy/go-idk/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	user := &domain.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, time.Now(), time.Now())

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Email, user.Name).
		WillReturnRows(rows)

	ctx := context.Background()
	err = repo.Create(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
		AddRow(1, "test@example.com", "Test User", now, now)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id").
		WithArgs(1).
		WillReturnRows(rows)

	ctx := context.Background()
	user, err := repo.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id").
		WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}))

	ctx := context.Background()
	user, err := repo.GetByID(ctx, 999)

	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
		AddRow(1, "test@example.com", "Test User", now, now)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE email").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	ctx := context.Background()
	user, err := repo.GetByEmail(ctx, "test@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT").
		WillReturnRows(countRows)

	// Mock select query
	now := time.Now()
	userRows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
		AddRow(1, "test1@example.com", "Test User 1", now, now).
		AddRow(2, "test2@example.com", "Test User 2", now, now)

	mock.ExpectQuery("SELECT (.+) FROM users ORDER BY id LIMIT").
		WithArgs(10, 0).
		WillReturnRows(userRows)

	ctx := context.Background()
	users, total, err := repo.GetAll(ctx, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, users, 2)
	assert.Equal(t, "test1@example.com", users[0].Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	user := &domain.User{
		Email: "updated@example.com",
		Name:  "Updated Name",
	}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"updated_at"}).AddRow(now)

	mock.ExpectQuery("UPDATE users").
		WithArgs(user.Email, user.Name, 1).
		WillReturnRows(rows)

	ctx := context.Background()
	err = repo.Update(ctx, 1, user)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	user := &domain.User{
		Email: "updated@example.com",
		Name:  "Updated Name",
	}

	mock.ExpectQuery("UPDATE users").
		WithArgs(user.Email, user.Name, 999).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}))

	ctx := context.Background()
	err = repo.Update(ctx, 999, user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	mock.ExpectExec("DELETE FROM users WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err = repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	mock.ExpectExec("DELETE FROM users WHERE id").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx := context.Background()
	err = repo.Delete(ctx, 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_WithContext_Timeout(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		errDB := db.Close()
		if errDB != nil {
			err = errDB
		}
	}()

	repo := NewUserRepository(&database.DB{DB: db})

	user := &domain.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Email, user.Name).
		WillReturnError(context.Canceled)

	err = repo.Create(ctx, user)

	assert.Error(t, err)
}
