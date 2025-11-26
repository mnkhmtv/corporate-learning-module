package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
)

type MentorHandler struct {
	mentorService *service.MentorService
}

func NewMentorHandler(mentorService *service.MentorService) *MentorHandler {
	return &MentorHandler{
		mentorService: mentorService,
	}
}

// CreateMentorDTO represents mentor creation input
type CreateMentorDTO struct {
	Name       string `json:"name" binding:"required"`
	JobTitle   string `json:"jobTitle" binding:"required"`
	Experience string `json:"experience"`
	Email      string `json:"email" binding:"required,email"`
	Telegram   string `json:"telegram"`
}

// GetAllMentors handles GET /api/mentors
func (h *MentorHandler) GetAllMentors(c *gin.Context) {
	// Check if filtering by availability
	availableOnly := c.Query("available") == "true"

	var mentors interface{}
	var err error

	if availableOnly {
		mentors, err = h.mentorService.GetAvailableMentors(c.Request.Context())
	} else {
		mentors, err = h.mentorService.GetAllMentors(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mentors)
}

// GetMentorByID handles GET /api/mentors/:id
func (h *MentorHandler) GetMentorByID(c *gin.Context) {
	mentorID := c.Param("id")

	mentor, err := h.mentorService.GetMentorByID(c.Request.Context(), mentorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mentor)
}

// CreateMentor handles POST /api/mentors (admin only)
func (h *MentorHandler) CreateMentor(c *gin.Context) {
	var req CreateMentorDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mentor, err := h.mentorService.CreateMentor(
		c.Request.Context(),
		req.Name, req.JobTitle, req.Experience,
		req.Email, req.Telegram,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mentor)
}
