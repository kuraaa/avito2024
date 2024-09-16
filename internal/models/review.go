// internal/models/review.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
