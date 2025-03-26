package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sparkfund/email-service/internal/models"
	"sparkfund/email-service/internal/services"
	"sparkfund/email-service/internal/errors"
)

type Handler struct {
	logger       *zap.Logger
	emailService services.EmailService
}

func NewHandler(logger *zap.Logger, emailService services.EmailService) *Handler {
	return &Handler{
		logger:       logger,
		emailService: emailService,
	}
}

func (h *Handler) handleError(c *gin.Context, err error) {
	var appError *errors.Error

	if errors.As(err, &appError) {
		h.logger.Error("Request failed",
			zap.Int("status", appError.Status),
			zap.String("message", appError.Message),
			zap.Error(err))
		c.JSON(appError.Status, models.Error{
			Code:    appError.Status,
			Message: appError.Message,
		})
		return
	}

	h.logger.Error("Request failed",
		zap.Int("status", http.StatusInternalServerError),
		zap.Error(err))
	c.JSON(http.StatusInternalServerError, models.Error{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	})
}

// SendEmail godoc
// @Summary      Send an email
// @Description  Send an email to one or more recipients
// @ID           sendEmail
// @Tags         emails
// @Accept       json
// @Produce      json
// @Param        email body models.SendEmailRequest true "Email request"
// @Success      200 {object} models.EmailResponse
// @Failure      400 {object} models.Error
// @Failure      500 {object} models.Error
// @Router       /emails [post]
// @Security     BearerAuth

// CreateTemplate godoc
// @Summary      Create an email template
// @Description  Create a new email template
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        template body models.CreateTemplateRequest true "Template request"
// @Success      201 {object} models.Template
// @Failure      400 {object} models.Error
// @Failure      500 {object} models.Error
// @Router       /templates [post]
// @Security     BearerAuth
func (h *Handler) SendEmail(c *gin.Context) {
	var req models.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
		})
		return
	}

	if err := models.ValidateStruct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := h.emailService.SendEmail(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	var req models.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
		})
		return
	}

	if err := models.ValidateStruct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	template, err := h.emailService.CreateTemplate(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, template)
}