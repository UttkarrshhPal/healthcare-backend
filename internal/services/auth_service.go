package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"healthcare-portal/internal/models"
	"healthcare-portal/internal/repository"
	"healthcare-portal/internal/utils"
)

type AuthService interface {
	Login(email, password string) (string, *models.User, error)
	Register(user *models.User) error
	ValidateToken(token string) (*utils.Claims, error)
	RefreshToken(userID uint) (string, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ResetPassword(email string) (string, error)
	VerifyResetToken(token string) (*models.User, error)
	UpdatePassword(userID uint, newPassword string) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(email, password string) (string, *models.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return "", nil, errors.New("user account is deactivated")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}

// Register creates a new user account
func (s *authService) Register(user *models.User) error {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("email already registered")
	}

	// Validate role
	if user.Role != models.RoleReceptionist && user.Role != models.RoleDoctor {
		return errors.New("invalid role")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	user.IsActive = true

	// Create user
	return s.userRepo.Create(user)
}

// ValidateToken validates a JWT token and returns the claims
func (s *authService) ValidateToken(token string) (*utils.Claims, error) {
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Check if user still exists and is active
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	return claims, nil
}

// RefreshToken generates a new JWT token for a user
func (s *authService) RefreshToken(userID uint) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.IsActive {
		return "", errors.New("user account is deactivated")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

// ChangePassword changes a user's password
func (s *authService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("incorrect current password")
	}

	// Validate new password
	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update password
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

// ResetPassword initiates password reset process
func (s *authService) ResetPassword(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return "", nil
	}

	if !user.IsActive {
		return "", nil
	}

	// Generate a temporary reset token (valid for 1 hour)
	resetToken, err := utils.GeneratePasswordResetToken(user.ID, user.Email)
	if err != nil {
		return "", errors.New("failed to generate reset token")
	}

	// In a real application, you would send this token via email
	// For now, we'll return it (in production, never return the token directly)
	return resetToken, nil
}

// VerifyResetToken verifies a password reset token
func (s *authService) VerifyResetToken(token string) (*models.User, error) {
	claims, err := utils.ValidatePasswordResetToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired reset token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	return user, nil
}

// UpdatePassword updates a user's password without requiring the old password
func (s *authService) UpdatePassword(userID uint, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Validate new password
	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update password
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

// GetUserByID retrieves a user by ID
func (s *authService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *authService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	return user, nil
}

// Add your JWT utility functions here, e.g., GenerateJWT, ValidateJWT, etc.

type PasswordResetClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// ValidatePasswordResetToken validates a password reset JWT token and returns the claims
func ValidatePasswordResetToken(tokenString string) (*PasswordResetClaims, error) {
	secret := []byte("your-reset-secret") // Use a secure secret and store in env/config

	token, err := jwt.ParseWithClaims(tokenString, &PasswordResetClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*PasswordResetClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired reset token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("reset token expired")
	}

	return claims, nil
}
