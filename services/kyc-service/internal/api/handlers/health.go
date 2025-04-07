package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type HealthResponse struct {
    Status    string `json:"status"`
    Version   string `json:"version"`
    CommitSHA string `json:"commitSha"`
}

func HealthCheck(version, commitSHA string) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, HealthResponse{
            Status:    "healthy",
            Version:   version,
            CommitSHA: commitSHA,
        })
    }
}