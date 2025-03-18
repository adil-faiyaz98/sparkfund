package handlers

import (
	"fmt"
	"net/http"

	"github.com/adil-faiyaz98/structgen/internal/reports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportHandler struct {
	service reports.ReportService
}

func NewReportHandler(service reports.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) CreateReport(c *gin.Context) {
	var report reports.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateReport(c.Request.Context(), &report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

func (h *ReportHandler) GetReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
		return
	}

	report, err := h.service.GetReport(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) GetUserReports(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	reports, err := h.service.GetUserReports(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (h *ReportHandler) UpdateReportStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
		return
	}

	var req struct {
		Status  string `json:"status" binding:"required"`
		FileURL string `json:"file_url"`
		Error   string `json:"error"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := reports.ReportStatus(req.Status)
	if !isValidReportStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report status"})
		return
	}

	var reportErr error
	if req.Error != "" {
		reportErr = fmt.Errorf(req.Error)
	}

	if err := h.service.UpdateReportStatus(c.Request.Context(), id, status, req.FileURL, reportErr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "report status updated successfully"})
}

func (h *ReportHandler) DeleteReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
		return
	}

	if err := h.service.DeleteReport(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "report deleted successfully"})
}

func isValidReportStatus(status reports.ReportStatus) bool {
	switch status {
	case reports.ReportStatusPending,
		reports.ReportStatusGenerating,
		reports.ReportStatusCompleted,
		reports.ReportStatusFailed:
		return true
	default:
		return false
	}
}
