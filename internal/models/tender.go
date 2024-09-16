// internal/models/tender.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Tender struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ServiceType    string    `json:"serviceType"`
	Status         string    `json:"status"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Version        int       `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatorUsername string   `json:"creatorUsername"`
}