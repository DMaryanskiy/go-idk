package domain

import (
	"context"
	"time"
)

// Entity
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"string"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DTOs (Data Transfer Object)
type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
	Name  string `json:"name" validate:"required,min=2,max=255"`
}

type UpdateUserRequest struct {
	Email string `json:"email" validate:"omitempty,email,max=255"`
	Name  string `json:"name" validate:"omitempty,min=2,max=255"`
}

type PaginationResponse struct {
	Users      []User `json:"users"`
	Total      int    `json:"total"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	TotalPages int    `json:"total_pages"`
}

// Repository interface (contract)
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context, limit, offset int) ([]User, int, error)
	Update(ctx context.Context, id int, user *User) error
	Delete(ctx context.Context, id int) error
}

// Service interface (contract)
type UserService interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetUser(ctx context.Context, id int) (*User, error)
	GetUsers(ctx context.Context, limit, offset int) (*PaginationResponse, error)
	UpdateUser(ctx context.Context, id int, req *UpdateUserRequest) (*User, error)
	DeleteUser(ctx context.Context, id int) error
}
