package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
	*BaseHandler
}

func NewMetricsHandler(base *BaseHandler) *MetricsHandler {
	return &MetricsHandler{BaseHandler: base}
}

func (h *MetricsHandler) RegisterRoutes(r *gin.RouterGroup) {
	metrics := r.Group("/metrics")
	{
		metrics.GET("", gin.WrapH(promhttp.Handler()))
		metrics.GET("/kyc", h.KYCMetrics)
	}
}
