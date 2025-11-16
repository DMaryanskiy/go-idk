package service

import (
	"fmt"
	"strings"

	"github.com/DMaryanskiy/go-idk/internal/domain"
	"go.uber.org/zap"
)

type userService struct {
	repo   domain.UserRepository
	logger *zap.Logger
}

func NewUserService(repo domain.UserRepository, logger *zap.Logger) domain.UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) CreateUser(req *domain.CreateUserRequest) (*domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.repo.GetByEmail(email)
	if err != nil {
		s.logger.Error("Error checking existing user", zap.Error(err))
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	user := &domain.User{
		Email: email,
		Name:  strings.TrimSpace(req.Name),
	}

	if err := s.repo.Create(user); err != nil {
		s.logger.Error("Error creating user", zap.Error(err))
		return nil, fmt.Errorf("failed to create a user: %w", err)
	}

	s.logger.Info("User created", zap.Int("user_id", user.ID), zap.String("email", user.Email))
	return user, nil
}

func (s *userService) GetUser(id int) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *userService) GetUsers(limit, offset int) (*domain.PaginationResponse, error) {
	if limit < 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.repo.GetAll(limit, offset)
	if err != nil {
		s.logger.Error("Error getting all users", zap.Error(err))
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	totalPages := (total + limit - 1) / limit

	return &domain.PaginationResponse{
		Users:      users,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) UpdateUser(id int, req *domain.UpdateUserRequest) (*domain.User, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Error getting user by id", zap.Error(err))
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("user not found")
	}

	if req.Email != "" {
		email := strings.ToLower(strings.TrimSpace(req.Email))

		if email != existing.Email {
			emailExists, err := s.repo.GetByEmail(email)
			if err != nil {
				s.logger.Error("Error checking email availability", zap.Error(err))
				return nil, fmt.Errorf("failed to check email availability: %w", err)
			}
			if emailExists != nil {
				return nil, fmt.Errorf("email %s already in use", email)
			}
			existing.Email = email
		}
	}
	if req.Name != "" {
		existing.Name = strings.TrimSpace(req.Name)
	}

	if err := s.repo.Update(id, existing); err != nil {
		s.logger.Error("Error updating user", zap.Int("user_id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to update user with id %d: %w", id, err)
	}

	s.logger.Info("User updated", zap.Int("user_id", id))
	return existing, nil
}

func (s *userService) DeleteUser(id int) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("Error deleting user", zap.Int("user_id", id), zap.Error(err))
		return fmt.Errorf("failed to delete user with id %d: %w", id, err)
	}

	s.logger.Info("User deleted", zap.Int("user_id", id))
	return nil
}
