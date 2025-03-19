package handlers

import (
	"fmt"
	"net/http"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Email Service API
// @version         1.0
// @description     A service for sending emails and managing email templates
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

type Handler struct {
	service services.EmailService
}

func NewHandler(service services.EmailService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		api.POST("/send", h.SendEmail)
		api.POST("/template", h.CreateTemplate)
		api.GET("/template/:id", h.GetTemplate)
		api.PUT("/template/:id", h.UpdateTemplate)
		api.DELETE("/template/:id", h.DeleteTemplate)
		api.GET("/logs", h.GetEmailLogs)
	}
}

// @Summary      Send an email
// @Description  Send an email with optional attachments
// @Tags         email
// @Accept       json
// @Produce      json
// @Param        request body models.SendEmailRequest true "Email request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /send [post]
func (h *Handler) SendEmail(c *gin.Context) {
	var req models.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": validationErr.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.service.SendEmail(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error sending email: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

// @Summary      Create an email template
// @Description  Create a new email template
// @Tags         template
// @Accept       json
// @Produce      json
// @Param        request body models.CreateTemplateRequest true "Template request"
// @Success      201  {object}  models.Template
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /template [post]
func (h *Handler) CreateTemplate(c *gin.Context) {
	var req models.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": validationErr.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	template, err := h.service.CreateTemplate(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating template: %v", err)})
		return
	}
	c.JSON(http.StatusCreated, template)
}

// @Summary      Get an email template
// @Description  Get an email template by ID
// @Tags         template
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Template ID"
// @Success      200  {object}  models.Template
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /template/{id} [get]
func (h *Handler) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing template ID"})
		return
	}

	template, err := h.service.GetTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting template: %v", err)})
		return
	}
	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}

// @Summary      Update an email template
// @Description  Update an existing email template
// @Tags         template
// @Accept       json
// @Produce      json
// @Param        id      path      string                    true  "Template ID"
// @Param        request body      models.UpdateTemplateRequest  true  "Template update request"
// @Success      200  {object}  models.Template
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /template/{id} [put]
func (h *Handler) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing template ID"})
		return
	}

	var req models.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": validationErr.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	template, err := h.service.UpdateTemplate(templateID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating template: %v", err)})
		return
	}
	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}

// @Summary      Delete an email template
// @Description  Delete an email template by ID
// @Tags         template
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Template ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /template/{id} [delete]
func (h *Handler) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing template ID"})
		return
	}

	err := h.service.DeleteTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error deleting template: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// @Summary      Get email logs
// @Description  Get all email logs
// @Tags         logs
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.EmailLog
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /logs [get]
func (h *Handler) GetEmailLogs(c *gin.Context) {
	logs, err := h.service.GetEmailLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting email logs: %v", err)})
		return
	}
	c.JSON(http.StatusOK, logs)
}
