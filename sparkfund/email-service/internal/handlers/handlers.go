package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/email-service/internal/models"
	"github.com/sparkfund/email-service/internal/services"
	"go.uber.org/zap"
)

// Handler handles HTTP requests
type Handler struct {
	logger       *zap.Logger
	emailService services.EmailService
}

// NewHandler creates a new handler instance
func NewHandler(logger *zap.Logger, emailService services.EmailService) *Handler {
	return &Handler{
		logger:       logger,
		emailService: emailService,
	}
}

// SendEmail godoc
// @Summary      Send an email
// @Description  Send an email to one or more recipients
// @Tags         emails
// @Accept       json
// @Produce      json
// @Param        email body models.SendEmailRequest true "Email request"
// @Success      200 {object} models.EmailResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /emails [post]
// @Security     BearerAuth
func (h *Handler) SendEmail(c *gin.Context) {
	var req models.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	if err := h.emailService.SendEmail(req); err != nil {
		h.logger.Error("Failed to send email", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to send email",
		})
		return
	}

	c.JSON(http.StatusOK, models.EmailResponse{
		Message: "Email queued for sending",
	})
}

// GetEmailLogs godoc
// @Summary      Get email logs
// @Description  Retrieve all email logs with optional filtering
// @Tags         emails
// @Accept       json
// @Produce      json
// @Param        status query string false "Filter by status"
// @Param        from query string false "Filter by start date"
// @Param        to query string false "Filter by end date"
// @Param        limit query int false "Limit number of results"
// @Param        offset query int false "Offset for pagination"
// @Success      200 {array} models.EmailLog
// @Failure      500 {object} models.ErrorResponse
// @Router       /emails [get]
// @Security     BearerAuth
func (h *Handler) GetEmailLogs(c *gin.Context) {
	logs, err := h.emailService.GetEmailLogs()
	if err != nil {
		h.logger.Error("Failed to get email logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get email logs",
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// CreateTemplate godoc
// @Summary      Create a new template
// @Description  Create a new email template
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        template body models.CreateTemplateRequest true "Template request"
// @Success      201 {object} models.Template
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /templates [post]
// @Security     BearerAuth
func (h *Handler) CreateTemplate(c *gin.Context) {
	var req models.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	template, err := h.emailService.CreateTemplate(req)
	if err != nil {
		h.logger.Error("Failed to create template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create template",
		})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// GetTemplate godoc
// @Summary      Get a template
// @Description  Retrieve a template by ID
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        id path string true "Template ID"
// @Success      200 {object} models.Template
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /templates/{id} [get]
// @Security     BearerAuth
func (h *Handler) GetTemplate(c *gin.Context) {
	id := c.Param("id")
	template, err := h.emailService.GetTemplate(id)
	if err != nil {
		h.logger.Error("Failed to get template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get template",
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// UpdateTemplate godoc
// @Summary      Update a template
// @Description  Update an existing template
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        id path string true "Template ID"
// @Param        template body models.UpdateTemplateRequest true "Template request"
// @Success      200 {object} models.Template
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /templates/{id} [put]
// @Security     BearerAuth
func (h *Handler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	template, err := h.emailService.UpdateTemplate(id, req)
	if err != nil {
		h.logger.Error("Failed to update template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update template",
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate godoc
// @Summary      Delete a template
// @Description  Delete a template by ID
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        id path string true "Template ID"
// @Success      204 "No Content"
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /templates/{id} [delete]
// @Security     BearerAuth
func (h *Handler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.emailService.DeleteTemplate(id); err != nil {
		h.logger.Error("Failed to delete template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete template",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
