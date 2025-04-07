package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/domain/model"
	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/service"
)

// ExampleHandler handles HTTP requests for examples
type ExampleHandler struct {
	service service.ExampleService
	logger  *logrus.Logger
}

// NewExampleHandler creates a new example handler
func NewExampleHandler(service service.ExampleService, logger *logrus.Logger) *ExampleHandler {
	return &ExampleHandler{
		service: service,
		logger:  logger,
	}
}

// GetAll returns all examples
func (h *ExampleHandler) GetAll(c *gin.Context) {
	examples, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get examples")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get examples"})
		return
	}

	c.JSON(http.StatusOK, examples)
}

// GetByID returns an example by ID
func (h *ExampleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	example, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get example")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get example"})
		return
	}

	if example == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Example not found"})
		return
	}

	c.JSON(http.StatusOK, example)
}

// Create creates a new example
func (h *ExampleHandler) Create(c *gin.Context) {
	var example model.Example
	if err := c.ShouldBindJSON(&example); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.Create(c.Request.Context(), &example); err != nil {
		h.logger.WithError(err).Error("Failed to create example")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create example"})
		return
	}

	c.JSON(http.StatusCreated, example)
}

// Update updates an example
func (h *ExampleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var example model.Example
	if err := c.ShouldBindJSON(&example); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	example.ID = id

	if err := h.service.Update(c.Request.Context(), &example); err != nil {
		h.logger.WithError(err).Error("Failed to update example")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update example"})
		return
	}

	c.JSON(http.StatusOK, example)
}

// Delete deletes an example
func (h *ExampleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.logger.WithError(err).Error("Failed to delete example")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete example"})
		return
	}

	c.Status(http.StatusNoContent)
}
