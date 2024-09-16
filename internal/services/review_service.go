// internal/services/review_service.go
package services

import (
	"avito-tender-service/internal/models"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ReviewService struct {
	db *sql.DB
}

func NewReviewService(db *sql.DB) *ReviewService {
	return &ReviewService{db: db}
}

func (s *ReviewService) CreateReview(review *models.Review) error {
	review.ID = uuid.New()
	review.CreatedAt = time.Now()

	_, err := s.db.Exec("INSERT INTO reviews (id, description, created_at) VALUES ($1, $2, $3)",
		review.ID, review.Description, review.CreatedAt)
	return err
}
