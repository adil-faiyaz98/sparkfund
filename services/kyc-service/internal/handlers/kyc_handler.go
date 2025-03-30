package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"kyc-service/internal/models"
	"kyc-service/internal/service"

	"github.com/go-chi/chi/v5"
)

// KYCHandler handles HTTP requests for KYC operations
type KYCHandler struct {
	service *service.KYCService
}

// NewKYCHandler creates a new KYC handler instance
func NewKYCHandler() *KYCHandler {
	return &KYCHandler{
		service: service.NewKYCService(),
	}
}

// CreateKYC handles the creation of a new KYC record
func (h *KYCHandler) CreateKYC(w http.ResponseWriter, r *http.Request) {
	var kyc models.KYC
	if err := json.NewDecoder(r.Body).Decode(&kyc); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateKYC(r.Context(), &kyc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(kyc)
}

// GetKYC handles retrieving a KYC record by user ID
func (h *KYCHandler) GetKYC(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	kyc, err := h.service.GetKYC(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if kyc == nil {
		http.Error(w, "KYC record not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(kyc)
}

// UpdateKYC handles updating an existing KYC record
func (h *KYCHandler) UpdateKYC(w http.ResponseWriter, r *http.Request) {
	var kyc models.KYC
	if err := json.NewDecoder(r.Body).Decode(&kyc); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateKYC(r.Context(), &kyc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(kyc)
}

// UpdateKYCStatus handles updating the status of a KYC record
func (h *KYCHandler) UpdateKYCStatus(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Status          string `json:"status"`
		RejectionReason string `json:"rejection_reason"`
		ReviewerID      string `json:"reviewer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateKYCStatus(r.Context(), userID, req.Status, req.RejectionReason, req.ReviewerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ListKYC handles retrieving a list of KYC records
func (h *KYCHandler) ListKYC(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	kycs, total, err := h.service.ListKYC(r.Context(), status, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Data  []models.KYC `json:"data"`
		Total int64        `json:"total"`
	}{
		Data:  kycs,
		Total: total,
	}

	json.NewEncoder(w).Encode(response)
}

// DeleteKYC handles deleting a KYC record
func (h *KYCHandler) DeleteKYC(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteKYC(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
} 