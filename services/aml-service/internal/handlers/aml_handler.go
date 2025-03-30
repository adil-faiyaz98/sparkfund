package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"aml-service/internal/models"
	"aml-service/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AMLHandler struct {
	service services.AMLService
	logger  *zap.Logger
}

func NewAMLHandler(service services.AMLService, logger *zap.Logger) *AMLHandler {
	return &AMLHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AMLHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/transactions", h.ProcessTransaction)
	r.Get("/transactions/{id}", h.GetTransaction)
	r.Get("/transactions", h.ListTransactions)
	r.Get("/risk-profiles/{userId}", h.GetRiskProfile)

	return r
}

func (h *AMLHandler) ProcessTransaction(w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ProcessTransaction(r.Context(), &tx); err != nil {
		h.logger.Error("Failed to process transaction", zap.Error(err))
		http.Error(w, "Failed to process transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

func (h *AMLHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("Invalid transaction ID", zap.Error(err))
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	tx, err := h.service.GetTransaction(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get transaction", zap.Error(err))
		http.Error(w, "Failed to get transaction", http.StatusInternalServerError)
		return
	}

	if tx == nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func (h *AMLHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Error("Invalid limit parameter", zap.Error(err))
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			h.logger.Error("Invalid offset parameter", zap.Error(err))
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	transactions, err := h.service.ListTransactions(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list transactions", zap.Error(err))
		http.Error(w, "Failed to list transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *AMLHandler) GetRiskProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	rp, err := h.service.GetRiskProfile(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get risk profile", zap.Error(err))
		http.Error(w, "Failed to get risk profile", http.StatusInternalServerError)
		return
	}

	if rp == nil {
		http.Error(w, "Risk profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rp)
}
