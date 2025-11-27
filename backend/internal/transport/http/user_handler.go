package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/dto"
)

type UserHandler struct {
	userService     *service.UserService
	learningService *service.LearningService
	requestService  *service.RequestService
}

func NewUserHandler(userService *service.UserService, learningService *service.LearningService, requestService *service.RequestService) *UserHandler {
	return &UserHandler{
		userService:     userService,
		learningService: learningService,
		requestService:  requestService,
	}
}

// GetAllUsers handles GET /api/users (admin only)
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUserByID handles GET /api/users/:id (admin only)
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserByID handles PUT /api/users/:id (admin only)
func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(
		c.Request.Context(),
		id,
		req.Name,
		req.Email,
		req.Department,
		req.JobTitle,
		req.Telegram,
		req.Password,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserRequests handles GET /api/users/:id/requests (admin only)
func (h *UserHandler) GetUserRequests(c *gin.Context) {
	userID := c.Param("id")

	requests, err := h.requestService.GetUserRequests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseDTOs := dto.ToRequestResponseDTOs(requests)
	c.JSON(http.StatusOK, gin.H{"requests": responseDTOs})
}

// GetUserLearnings handles GET /api/users/:id/learnings (admin only)
func (h *UserHandler) GetUserLearnings(c *gin.Context) {
	userID := c.Param("id")

	learnings, err := h.learningService.GetUserLearnings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseDTOs := dto.ToLearningResponseDTOs(learnings)
	c.JSON(http.StatusOK, gin.H{"learnings": responseDTOs})
}
