package handlers

import (
	"net/http"

	"healthcare-portal/internal/models"
	"healthcare-portal/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// @Summary Login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  *user,
	})
}

type RegisterRequest struct {
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=6"`
	Name     string          `json:"name" binding:"required"`
	Role     models.UserRole `json:"role" binding:"required"`
}

// @Summary Register
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} models.User
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
	}

	if err := h.authService.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = "" // Remove password from response
	c.JSON(http.StatusCreated, user)
}

// @Summary Get Current User
// @Description Get currently logged in user details
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Router /api/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	user := models.User{
		ID:    userID.(uint),
		Email: email.(string),
		Role:  models.UserRole(role.(string)),
	}

	c.JSON(http.StatusOK, user)
}
