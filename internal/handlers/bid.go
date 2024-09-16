package handlers

import (
	"avito-tender-service/internal/models"
	"avito-tender-service/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type BidHandler struct {
	service *services.BidService
}

func NewBidHandler(service *services.BidService) *BidHandler {
	return &BidHandler{service: service}
}

func (h *BidHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	var bid models.Bid
	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateBid(&bid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bid)
}

func (h *BidHandler) GetUserBids(w http.ResponseWriter, r *http.Request) {
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

	bids, err := h.service.GetUserBids(username, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bids)
}

func (h *BidHandler) GetBidsForTender(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, "tenderId")
	if tenderId == "" {
		http.Error(w, "tenderId is required", http.StatusBadRequest)
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

	bids, err := h.service.GetBidsForTender(tenderId, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bids)
}

func (h *BidHandler) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	bidId := chi.URLParam(r, "bidId")
	if bidId == "" {
		http.Error(w, "bidId is required", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetBidStatus(bidId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (h *BidHandler) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	bidId := chi.URLParam(r, "bidId")
	if bidId == "" {
		http.Error(w, "bidId is required", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateBidStatus(bidId, status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (h *BidHandler) EditBid(w http.ResponseWriter, r *http.Request) {
	bidId := chi.URLParam(r, "bidId")
	if bidId == "" {
		http.Error(w, "bidId is required", http.StatusBadRequest)
		return
	}

	var bid models.Bid
	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bid.ID = uuid.MustParse(bidId)
	bid.UpdatedAt = time.Now()
	bid.Version++

	if err := h.service.EditBid(&bid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func (h *BidHandler) SubmitBidDecision(w http.ResponseWriter, r *http.Request) {
	bidId := chi.URLParam(r, "bidId")
	if bidId == "" {
		http.Error(w, "bidId is required", http.StatusBadRequest)
		return
	}

	decision := r.URL.Query().Get("decision")
	if decision == "" {
		http.Error(w, "decision is required", http.StatusBadRequest)
		return
	}

	if err := h.service.SubmitBidDecision(bidId, decision); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"decision": decision})
}

func (h *BidHandler) SubmitBidFeedback(w http.ResponseWriter, r *http.Request) {
	bidId := chi.URLParam(r, "bidId")
	if bidId == "" {
		http.Error(w, "bidId is required", http.StatusBadRequest)
		return
	}

	feedback := r.URL.Query().Get("feedback")
	if feedback == "" {
		http.Error(w, "feedback is required", http.StatusBadRequest)
		return
	}

	if err := h.service.SubmitBidFeedback(bidId, feedback); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"feedback": feedback})
}

func (h *BidHandler) RollbackBid(w http.ResponseWriter, r *http.Request) {
	bidId := chi.URLParam(r, "bidId")
	if bidId == "" {
		http.Error(w, "bidId is required", http.StatusBadRequest)
		return
	}

	versionStr := chi.URLParam(r, "version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "invalid version", http.StatusBadRequest)
		return
	}

	bid, err := h.service.RollbackBid(bidId, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func (h *BidHandler) GetBidReviews(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, "tenderId")
	if tenderId == "" {
		http.Error(w, "tenderId is required", http.StatusBadRequest)
		return
	}

	authorUsername := r.URL.Query().Get("authorUsername")
	if authorUsername == "" {
		http.Error(w, "authorUsername is required", http.StatusBadRequest)
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

	reviews, err := h.service.GetBidReviews(tenderId, authorUsername, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviews)
}