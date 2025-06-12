package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"healthcare-portal/internal/models"
	"healthcare-portal/internal/services"
)

// Mock repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByRole(role models.UserRole) ([]models.User, error) {
	args := m.Called(role)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAuthService(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	t.Run("Login Success", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		mockUser := &models.User{
			ID:       1,
			Email:    "test@example.com",
			Password: string(hashedPassword),
			Role:     models.RoleReceptionist,
			IsActive: true,
		}

		mockRepo.On("FindByEmail", "test@example.com").Return(mockUser, nil)

		token, user, err := authService.Login("test@example.com", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, mockUser.Email, user.Email)
	})

	t.Run("Register Success", func(t *testing.T) {
		newUser := &models.User{
			Email:    "new@example.com",
			Password: "password123",
			Name:     "New User",
			Role:     models.RoleDoctor,
		}

		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		err := authService.Register(newUser)
		assert.NoError(t, err)
		assert.NotEqual(t, "password123", newUser.Password) // Password should be hashed
	})
}
