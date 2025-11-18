package service

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/DMaryanskiy/go-idk/internal/domain"
    "go.uber.org/zap"
)

// Mock Repository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
    args := m.Called(ctx, limit, offset)
    return args.Get(0).([]domain.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, id int, user *domain.User) error {
    args := m.Called(ctx, id, user)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func TestCreateUser_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    req := &domain.CreateUserRequest{
        Email: "test@example.com",
        Name:  "Test User",
    }

    ctx := context.Background()

    mockRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, nil)
    mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).
        Run(func(args mock.Arguments) {
            user := args.Get(1).(*domain.User)
            user.ID = 1
            user.CreatedAt = time.Now()
            user.UpdatedAt = time.Now()
        }).
        Return(nil)

    user, err := service.CreateUser(ctx, req)
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
    assert.Equal(t, "Test User", user.Name)
    assert.Equal(t, 1, user.ID)
    mockRepo.AssertExpectations(t)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    req := &domain.CreateUserRequest{
        Email: "test@example.com",
        Name:  "Test User",
    }

    ctx := context.Background()
    existingUser := &domain.User{ID: 1, Email: "test@example.com"}
    
    mockRepo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)

    user, err := service.CreateUser(ctx, req)
    
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "already exists")
    mockRepo.AssertExpectations(t)
}

func TestCreateUser_RepositoryError(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    req := &domain.CreateUserRequest{
        Email: "test@example.com",
        Name:  "Test User",
    }

    ctx := context.Background()
    
    mockRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, nil)
    mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).
        Return(errors.New("database error"))

    user, err := service.CreateUser(ctx, req)
    
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "failed to create a user")
    mockRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    expectedUser := &domain.User{
        ID:        1,
        Email:     "test@example.com",
        Name:      "Test User",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    ctx := context.Background()
    mockRepo.On("GetByID", ctx, 1).Return(expectedUser, nil)

    user, err := service.GetUser(ctx, 1)
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, expectedUser.Email, user.Email)
    assert.Equal(t, expectedUser.Name, user.Name)
    mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    ctx := context.Background()
    mockRepo.On("GetByID", ctx, 999).Return(nil, nil)

    user, err := service.GetUser(ctx, 999)
    
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "not found")
    mockRepo.AssertExpectations(t)
}

func TestGetUser_RepositoryError(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    ctx := context.Background()
    mockRepo.On("GetByID", ctx, 1).Return(nil, errors.New("database error"))

    user, err := service.GetUser(ctx, 1)
    
    assert.Error(t, err)
    assert.Nil(t, user)
    mockRepo.AssertExpectations(t)
}

func TestGetUsers_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    expectedUsers := []domain.User{
        {ID: 1, Email: "test1@example.com", Name: "User 1"},
        {ID: 2, Email: "test2@example.com", Name: "User 2"},
    }

    ctx := context.Background()
    mockRepo.On("GetAll", ctx, 10, 0).Return(expectedUsers, 2, nil)

    response, err := service.GetUsers(ctx, 10, 0)
    
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, 2, response.Total)
    assert.Len(t, response.Users, 2)
    assert.Equal(t, 10, response.Limit)
    assert.Equal(t, 0, response.Offset)
    assert.Equal(t, 1, response.TotalPages)
    mockRepo.AssertExpectations(t)
}

func TestGetUsers_WithPagination(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    expectedUsers := []domain.User{
        {ID: 11, Email: "test11@example.com", Name: "User 11"},
    }

    ctx := context.Background()
    mockRepo.On("GetAll", ctx, 10, 10).Return(expectedUsers, 25, nil)

    response, err := service.GetUsers(ctx, 10, 10)
    
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, 25, response.Total)
    assert.Equal(t, 3, response.TotalPages) // 25 items / 10 per page = 3 pages
    mockRepo.AssertExpectations(t)
}

func TestGetUsers_InvalidLimit(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    ctx := context.Background()
    // Should default to limit=10
    mockRepo.On("GetAll", ctx, 10, 0).Return([]domain.User{}, 0, nil)

    response, err := service.GetUsers(ctx, 0, 0)
    
    assert.NoError(t, err)
    assert.Equal(t, 10, response.Limit)
    mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    existingUser := &domain.User{
        ID:    1,
        Email: "old@example.com",
        Name:  "Old Name",
    }

    req := &domain.UpdateUserRequest{
        Email: "new@example.com",
        Name:  "New Name",
    }

    ctx := context.Background()
    mockRepo.On("GetByID", ctx, 1).Return(existingUser, nil)
    mockRepo.On("GetByEmail", ctx, "new@example.com").Return(nil, nil)
    mockRepo.On("Update", ctx, 1, mock.AnythingOfType("*domain.User")).Return(nil)

    user, err := service.UpdateUser(ctx, 1, req)
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "new@example.com", user.Email)
    assert.Equal(t, "New Name", user.Name)
    mockRepo.AssertExpectations(t)
}

func TestUpdateUser_NotFound(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    req := &domain.UpdateUserRequest{
        Name: "New Name",
    }

    ctx := context.Background()
    mockRepo.On("GetByID", ctx, 999).Return(nil, nil)

    user, err := service.UpdateUser(ctx, 999, req)
    
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "not found")
    mockRepo.AssertExpectations(t)
}

func TestUpdateUser_EmailAlreadyInUse(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    existingUser := &domain.User{
        ID:    1,
        Email: "old@example.com",
        Name:  "Old Name",
    }

    anotherUser := &domain.User{
        ID:    2,
        Email: "taken@example.com",
        Name:  "Another User",
    }

    req := &domain.UpdateUserRequest{
        Email: "taken@example.com",
    }

    ctx := context.Background()
    mockRepo.On("GetByID", ctx, 1).Return(existingUser, nil)
    mockRepo.On("GetByEmail", ctx, "taken@example.com").Return(anotherUser, nil)

    user, err := service.UpdateUser(ctx, 1, req)
    
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "already in use")
    mockRepo.AssertExpectations(t)
}

func TestUpdateUser_PartialUpdate(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    existingUser := &domain.User{
        ID:    1,
        Email: "old@example.com",
        Name:  "Old Name",
    }

    req := &domain.UpdateUserRequest{
        Name: "New Name", // Only updating name
    }

    ctx := context.Background()
    mockRepo.On("GetByID", ctx, 1).Return(existingUser, nil)
    mockRepo.On("Update", ctx, 1, mock.AnythingOfType("*domain.User")).Return(nil)

    user, err := service.UpdateUser(ctx, 1, req)
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "old@example.com", user.Email) // Email unchanged
    assert.Equal(t, "New Name", user.Name)          // Name updated
    mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    ctx := context.Background()
    mockRepo.On("Delete", ctx, 1).Return(nil)

    err := service.DeleteUser(ctx, 1)
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    ctx := context.Background()
    mockRepo.On("Delete", ctx, 999).Return(errors.New("user not found"))

    err := service.DeleteUser(ctx, 999)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "failed to delete user")
    mockRepo.AssertExpectations(t)
}

func TestServiceWithContextCancellation(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel immediately

    mockRepo.On("GetByID", ctx, 1).Return(nil, context.Canceled)

    user, err := service.GetUser(ctx, 1)
    
    assert.Error(t, err)
    assert.Nil(t, user)
}

func TestServiceWithContextTimeout(t *testing.T) {
    mockRepo := new(MockUserRepository)
    logger, _ := zap.NewDevelopment()
    service := NewUserService(mockRepo, logger)

    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
    defer cancel()
    
    time.Sleep(2 * time.Nanosecond) // Ensure timeout

    mockRepo.On("GetByID", ctx, 1).Return(nil, context.DeadlineExceeded)

    user, err := service.GetUser(ctx, 1)
    
    assert.Error(t, err)
    assert.Nil(t, user)
}
