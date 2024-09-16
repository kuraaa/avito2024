// internal/services/tender_service.go
package services

import (
	"avito-tender-service/internal/models"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TenderService struct {
	db *sql.DB
}

func NewTenderService(db *sql.DB) *TenderService {
	return &TenderService{db: db}
}

func (s *TenderService) CreateTender(tender *models.Tender) error {
	tender.ID = uuid.New()
	tender.CreatedAt = time.Now()
	tender.UpdatedAt = time.Now()
	tender.Version = 1
	tender.Status = "Created"

	_, err := s.db.Exec("INSERT INTO tenders (id, name, description, service_type, status, organization_id, version, created_at, updated_at, creator_username) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		tender.ID, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.Version, tender.CreatedAt, tender.UpdatedAt, tender.CreatorUsername)
	return err
}

func (s *TenderService) GetTenders(limit, offset int, serviceTypes []string) ([]models.Tender, error) {
	query := "SELECT id, name, description, service_type, status, organization_id, version, created_at, updated_at, creator_username FROM tenders"
	var args []interface{}

	if len(serviceTypes) > 0 {
		placeholders := make([]string, len(serviceTypes))
		for i := range serviceTypes {
			placeholders[i] = "$" + strconv.Itoa(i+1)
			args = append(args, serviceTypes[i])
		}
		query += " WHERE service_type IN (" + strings.Join(placeholders, ",") + ")"
	}

	query += " ORDER BY name LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	log.Printf("Executing query: %s with args: %v", query, args)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.UpdatedAt, &tender.CreatorUsername)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (s *TenderService) GetUserTenders(username string, limit, offset int) ([]models.Tender, error) {
	query := "SELECT id, name, description, service_type, status, organization_id, version, created_at, updated_at, creator_username FROM tenders WHERE creator_username = $1 ORDER BY name LIMIT $2 OFFSET $3"
	rows, err := s.db.Query(query, username, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.UpdatedAt, &tender.CreatorUsername)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (s *TenderService) GetTenderStatus(tenderId string) (string, error) {
	var status string
	err := s.db.QueryRow("SELECT status FROM tenders WHERE id = $1", tenderId).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (s *TenderService) UpdateTenderStatus(tenderId, status string) error {
	_, err := s.db.Exec("UPDATE tenders SET status = $1 WHERE id = $2", status, tenderId)
	return err
}

func (s *TenderService) GetTenderByID(tenderId string) (*models.Tender, error) {
	var tender models.Tender
	err := s.db.QueryRow("SELECT id, name, description, service_type, status, organization_id, version, created_at, updated_at, creator_username FROM tenders WHERE id = $1", tenderId).
		Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.UpdatedAt, &tender.CreatorUsername)
	if err != nil {
		return nil, err
	}
	return &tender, nil
}

func (s *TenderService) UpdateTender(tender *models.Tender) error {
	_, err := s.db.Exec("UPDATE tenders SET name = $1, description = $2, service_type = $3, status = $4, organization_id = $5, version = $6, updated_at = $7 WHERE id = $8",
		tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.Version, time.Now(), tender.ID)
	return err
}

func (s *TenderService) DeleteTender(tenderId string) error {
	_, err := s.db.Exec("DELETE FROM tenders WHERE id = $1", tenderId)
	return err
}

func (s *TenderService) EditTender(tender *models.Tender) error {
	_, err := s.db.Exec("UPDATE tenders SET name = $1, description = $2, service_type = $3, status = $4, organization_id = $5, version = $6, updated_at = $7 WHERE id = $8",
		tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.Version+1, time.Now(), tender.ID)
	return err
}

func (s *TenderService) RollbackTender(tenderId string, version int) (*models.Tender, error) {
	var tender models.Tender
	err := s.db.QueryRow("SELECT id, name, description, service_type, status, organization_id, version, created_at, updated_at, creator_username FROM tenders WHERE id = $1 AND version = $2", tenderId, version).
		Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.UpdatedAt, &tender.CreatorUsername)
	if err != nil {
		return nil, err
	}

	tender.Version++
	tender.UpdatedAt = time.Now()

	_, err = s.db.Exec("UPDATE tenders SET name = $1, description = $2, service_type = $3, status = $4, organization_id = $5, version = $6, updated_at = $7 WHERE id = $8",
		tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.Version, tender.UpdatedAt, tender.ID)
	if err != nil {
		return nil, err
	}

	return &tender, nil
}