package domain

import "time"

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
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll(limit, offset int) ([]User, int, error)
	Update(id int, user *User) error
	Delete(id int) error
}

// Service interface (contract)
type UserService interface {
	CreateUser(req *CreateUserRequest) (*User, error)
	GetUser(id int) (*User, error)
	GetUsers(limit, offset int) (*PaginationResponse, error)
	UpdateUser(id int, req *UpdateUserRequest) (*User, error)
	DeleteUser(id int) error
}
