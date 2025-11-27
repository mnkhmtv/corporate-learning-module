package http

import (
	"net/http"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UpdateUserDTO represents user update input
type UpdateUserDTO struct {
	Name       *string `json:"name"`
	Email      *string `json:"email"`
	Password   *string `json:"password"`
	Department *string `json:"department"`
	JobTitle   *string `json:"jobTitle"`
	Telegram   *string `json:"telegram"`
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req UpdateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateCurrentUser(
		c.Request.Context(),
		userID.(string),
		req.Name, req.Email, req.Department, req.JobTitle, req.Telegram, req.Password,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(
		c.Request.Context(),
		id,
		req.Name, req.Email, req.Department, req.JobTitle, req.Telegram, req.Password,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
