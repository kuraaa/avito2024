// internal/handlers/tender.go
package handlers

import (
	"avito-tender-service/internal/models"
	"avito-tender-service/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type TenderHandler struct {
	service *services.TenderService
}

func NewTenderHandler(service *services.TenderService) *TenderHandler {
	return &TenderHandler{service: service}
}

func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if tender.CreatorUsername == "" {
		http.Error(w, "creatorUsername is required", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTender(&tender); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tender)
}

func (h *TenderHandler) GetTenders(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	serviceTypeStr := r.URL.Query().Get("service_type")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limitStr == "" {
		limit = 5 // Default value
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offsetStr == "" {
		offset = 0 // Default value
	}

	var serviceTypes []string
	if serviceTypeStr != "" {
		serviceTypes = strings.Split(serviceTypeStr, ",")
	}

	log.Printf("Getting tenders with limit: %d, offset: %d, serviceTypes: %v", limit, offset, serviceTypes)

	tenders, err := h.service.GetTenders(limit, offset, serviceTypes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)
}

func (h *TenderHandler) GetUserTenders(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limitStr == "" {
		limit = 5 // Default value
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offsetStr == "" {
		offset = 0 // Default value
	}

	log.Printf("Getting user tenders for username: %s with limit: %d, offset: %d", username, limit, offset)

	tenders, err := h.service.GetUserTenders(username, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)
}

func (h *TenderHandler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, "tenderId")
	if tenderId == "" {
		http.Error(w, "tenderId is required", http.StatusBadRequest)
		return
	}

	log.Printf("Getting tender status for tenderId: %s", tenderId)

	status, err := h.service.GetTenderStatus(tenderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (h *TenderHandler) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, "tenderId")
	status := r.URL.Query().Get("status")
	if tenderId == "" || status == "" {
		http.Error(w, "tenderId and status are required", http.StatusBadRequest)
		return
	}

	log.Printf("Updating tender status for tenderId: %s to status: %s", tenderId, status)

	if err := h.service.UpdateTenderStatus(tenderId, status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (h *TenderHandler) GetTenderByID(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, "tenderId")
	if tenderId == "" {
		http.Error(w, "tenderId is required", http.StatusBadRequest)
		return
	}

	tender, err := h.service.GetTenderByID(tenderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

func (h *TenderHandler) UpdateTender(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, "tenderId")
	if tenderId == "" {
		http.Error(w, "tenderId is required", http.StatusBadRequest)
		return
	}

	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tender.ID = uuid.MustParse(tenderId)
	tender.UpdatedAt = time.Now()

	if err := h.service.UpdateTender(&tender); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

func (h *TenderHandler) DeleteTender(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, "tenderId")
	if tenderId == "" {
		http.Error(w, "tenderId is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTender(tenderId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
func PingHandler(w http.ResponseWriter, r *http.Request) {
	if !isServerReady() {
		http.Error(w, "Server not ready", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func isServerReady() bool {
	// Функция, которая проверяет готовность сервера
	return true
}

func (h *TenderHandler) EditTender(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, "tenderId")
	if tenderId == "" {
		http.Error(w, "tenderId is required", http.StatusBadRequest)
		return
	}

	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tender.ID = uuid.MustParse(tenderId)
	tender.UpdatedAt = time.Now()
	tender.Version++

	if err := h.service.EditTender(&tender); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

func (h *TenderHandler) RollbackTender(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, "tenderId")
	if tenderId == "" {
		http.Error(w, "tenderId is required", http.StatusBadRequest)
		return
	}

	versionStr := chi.URLParam(r, "version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "invalid version", http.StatusBadRequest)
		return
	}

	tender, err := h.service.RollbackTender(tenderId, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}
