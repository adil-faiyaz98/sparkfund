package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	port            = 8081
	aiServiceURL    = "http://localhost:8001"
	aiServiceAPIKey = "your-api-key"
)

// DeviceInfo represents device information
type DeviceInfo struct {
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	DeviceType   string `json:"device_type,omitempty"`
	OS           string `json:"os,omitempty"`
	Browser      string `json:"browser,omitempty"`
	Location     string `json:"location,omitempty"`
	CapturedTime string `json:"captured_time,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	DeviceInfo DeviceInfo `json:"device_info"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	} `json:"user"`
}

// VerificationRequest represents a verification request
type VerificationRequest struct {
	UserID     string `json:"user_id"`
	KycID      string `json:"kyc_id"`
	DocumentID string `json:"document_id"`
	Method     string `json:"method"`
	Status     string `json:"status"`
}

// VerificationResponse represents a verification response
type VerificationResponse struct {
	Verification struct {
		ID         string `json:"id"`
		UserID     string `json:"user_id"`
		KycID      string `json:"kyc_id"`
		DocumentID string `json:"document_id"`
		Method     string `json:"method"`
		Status     string `json:"status"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	} `json:"verification"`
}

// DocumentAnalysisRequest represents a document analysis request
type DocumentAnalysisRequest struct {
	DocumentID     string `json:"document_id"`
	VerificationID string `json:"verification_id"`
}

// FaceMatchRequest represents a face match request
type FaceMatchRequest struct {
	DocumentID     string `json:"document_id"`
	SelfieID       string `json:"selfie_id"`
	VerificationID string `json:"verification_id"`
}

// RiskAnalysisRequest represents a risk analysis request
type RiskAnalysisRequest struct {
	UserID         string     `json:"user_id"`
	VerificationID string     `json:"verification_id"`
	DeviceInfo     DeviceInfo `json:"device_info"`
}

// AnomalyDetectionRequest represents an anomaly detection request
type AnomalyDetectionRequest struct {
	UserID         string     `json:"user_id"`
	VerificationID string     `json:"verification_id"`
	DeviceInfo     DeviceInfo `json:"device_info"`
}

// ProcessDocumentRequest represents a document processing request
type ProcessDocumentRequest struct {
	DocumentID     string     `json:"document_id"`
	SelfieID       string     `json:"selfie_id"`
	VerificationID string     `json:"verification_id"`
	DeviceInfo     DeviceInfo `json:"device_info"`
}

func main() {
	// Register API handlers
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/v1/auth/login", loginHandler)
	http.HandleFunc("/api/v1/ai/models", aiModelsHandler)
	http.HandleFunc("/api/v1/verifications", verificationsHandler)
	http.HandleFunc("/api/v1/ai/analyze-document", analyzeDocumentHandler)
	http.HandleFunc("/api/v1/ai/match-faces", matchFacesHandler)
	http.HandleFunc("/api/v1/ai/analyze-risk", analyzeRiskHandler)
	http.HandleFunc("/api/v1/ai/detect-anomalies", detectAnomaliesHandler)
	http.HandleFunc("/api/v1/ai/process-document", processDocumentHandler)
	http.HandleFunc("/api/v1/get-api-key", getApiKeyHandler)

	// Serve Swagger UI
	http.HandleFunc("/swagger-ui.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger-ui.html")
	})

	// Serve Swagger JSON
	http.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.json")
	})

	fmt.Printf("Starting KYC service on port %d...\n", port)
	fmt.Printf("Swagger UI available at http://localhost:%d/swagger-ui.html\n", port)
	fmt.Printf("API Key available at http://localhost:%d/api/v1/get-api-key\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"service": "kyc-service",
		"version": "1.0.0",
	}

	jsonResponse(w, response)
}

func getApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"api_key": aiServiceAPIKey,
		"note":    "Use this API key in the X-API-Key header when calling the AI service",
	}

	jsonResponse(w, response)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate login
	response := LoginResponse{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		User: struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Role      string `json:"role"`
		}{
			ID:        "123e4567-e89b-12d3-a456-426614174000",
			Email:     req.Email,
			FirstName: "John",
			LastName:  "Doe",
			Role:      "user",
		},
	}

	jsonResponse(w, response)
}

func aiModelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call AI service to get models
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/models", aiServiceURL), nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	// Add API key header
	req.Header.Add("X-API-Key", aiServiceAPIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling AI service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(body)), http.StatusInternalServerError)
		return
	}

	// Forward the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func verificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req VerificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Simulate verification creation
		response := VerificationResponse{
			Verification: struct {
				ID         string `json:"id"`
				UserID     string `json:"user_id"`
				KycID      string `json:"kyc_id"`
				DocumentID string `json:"document_id"`
				Method     string `json:"method"`
				Status     string `json:"status"`
				CreatedAt  string `json:"created_at"`
				UpdatedAt  string `json:"updated_at"`
			}{
				ID:         "123e4567-e89b-12d3-a456-426614174001",
				UserID:     req.UserID,
				KycID:      req.KycID,
				DocumentID: req.DocumentID,
				Method:     req.Method,
				Status:     req.Status,
				CreatedAt:  time.Now().Format(time.RFC3339),
				UpdatedAt:  time.Now().Format(time.RFC3339),
			},
		}

		jsonResponse(w, response)
	} else if r.Method == http.MethodGet {
		// Simulate verification retrieval
		response := VerificationResponse{
			Verification: struct {
				ID         string `json:"id"`
				UserID     string `json:"user_id"`
				KycID      string `json:"kyc_id"`
				DocumentID string `json:"document_id"`
				Method     string `json:"method"`
				Status     string `json:"status"`
				CreatedAt  string `json:"created_at"`
				UpdatedAt  string `json:"updated_at"`
			}{
				ID:         "123e4567-e89b-12d3-a456-426614174001",
				UserID:     "123e4567-e89b-12d3-a456-426614174000",
				KycID:      "123e4567-e89b-12d3-a456-426614174002",
				DocumentID: "123e4567-e89b-12d3-a456-426614174003",
				Method:     "AI",
				Status:     "PENDING",
				CreatedAt:  time.Now().Format(time.RFC3339),
				UpdatedAt:  time.Now().Format(time.RFC3339),
			},
		}

		jsonResponse(w, response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func analyzeDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DocumentAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create request to AI service
	reqBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling request: %v", err), http.StatusInternalServerError)
		return
	}

	// Create request
	aiReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/document/analyze-base64", aiServiceURL), bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	// Add headers
	aiReq.Header.Set("Content-Type", "application/json")
	aiReq.Header.Set("X-API-Key", aiServiceAPIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(aiReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling AI service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(body)), http.StatusInternalServerError)
		return
	}

	// Forward the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func matchFacesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FaceMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create request to AI service
	reqBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling request: %v", err), http.StatusInternalServerError)
		return
	}

	// Create request
	aiReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/face/match-base64", aiServiceURL), bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	// Add headers
	aiReq.Header.Set("Content-Type", "application/json")
	aiReq.Header.Set("X-API-Key", aiServiceAPIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(aiReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling AI service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(body)), http.StatusInternalServerError)
		return
	}

	// Forward the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func analyzeRiskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RiskAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create request to AI service
	reqBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling request: %v", err), http.StatusInternalServerError)
		return
	}

	// Create request
	aiReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/risk/analyze", aiServiceURL), bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	// Add headers
	aiReq.Header.Set("Content-Type", "application/json")
	aiReq.Header.Set("X-API-Key", aiServiceAPIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(aiReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling AI service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(body)), http.StatusInternalServerError)
		return
	}

	// Forward the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func detectAnomaliesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnomalyDetectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create request to AI service
	reqBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling request: %v", err), http.StatusInternalServerError)
		return
	}

	// Create request
	aiReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/anomaly/detect", aiServiceURL), bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	// Add headers
	aiReq.Header.Set("Content-Type", "application/json")
	aiReq.Header.Set("X-API-Key", aiServiceAPIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(aiReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling AI service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(body)), http.StatusInternalServerError)
		return
	}

	// Forward the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func processDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ProcessDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create mock response
	response := map[string]interface{}{
		"id":           req.VerificationID,
		"status":       "COMPLETED",
		"notes":        "All verification checks passed",
		"completed_at": time.Now().Format(time.RFC3339),
	}

	jsonResponse(w, response)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
