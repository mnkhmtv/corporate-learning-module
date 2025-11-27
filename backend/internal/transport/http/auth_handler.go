package http

import (
	"net/http"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(
		c.Request.Context(),
		req.Name, req.Email, req.Password, req.Role,
		req.Department, req.JobTitle, req.Telegram,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetMe handles GET /api/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
