// internal/handlers/review.go
package handlers

import (
	"avito-tender-service/internal/models"
	"avito-tender-service/internal/services"
	"encoding/json"
	"net/http"
)

type ReviewHandler struct {
	service *services.ReviewService
}

func NewReviewHandler(service *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateReview(&review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}
