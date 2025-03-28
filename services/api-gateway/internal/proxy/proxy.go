package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var investmentServiceURL = os.Getenv("INVESTMENT_SERVICE_URL")
var investmentServicePort = os.Getenv("INVESTMENT_SERVICE_PORT")

func init() {
	if investmentServiceURL == "" {
		investmentServiceURL = "http://investment-service"
	}
	if investmentServicePort == "" {
		investmentServicePort = "8080"
	}
}

// ProxyToInvestmentService forwards the request to the Investment Service
func ProxyToInvestmentService(c *gin.Context) {
	// Construct the target URL
	targetURL := fmt.Sprintf("%s:%s%s", investmentServiceURL, investmentServicePort, c.Request.URL.Path)
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Create a new request
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers from the original request
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set content type if not already set
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Create HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set response status code
	c.Status(resp.StatusCode)

	// Write response body
	if len(body) > 0 {
		// Try to pretty print JSON responses
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "", "  "); err == nil {
			c.Data(resp.StatusCode, "application/json", prettyJSON.Bytes())
		} else {
			c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
		}
	}
} 