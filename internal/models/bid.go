package models

import (
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TenderID    uuid.UUID `json:"tender_id"`
	AuthorType  string    `json:"author_type"`
	AuthorID    uuid.UUID `json:"author_id"`
	Status      string    `json:"status"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BidReview struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}