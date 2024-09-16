package services

import (
	"avito-tender-service/internal/models"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type BidService struct {
	db *sql.DB
}

func NewBidService(db *sql.DB) *BidService {
	return &BidService{db: db}
}

func (s *BidService) CreateBid(bid *models.Bid) error {
	bid.ID = uuid.New()
	bid.CreatedAt = time.Now()
	bid.UpdatedAt = time.Now()
	bid.Version = 1
	bid.Status = "Created"

	_, err := s.db.Exec("INSERT INTO bids (id, name, description, tender_id, author_type, author_id, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		bid.ID, bid.Name, bid.Description, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.Status, bid.Version, bid.CreatedAt, bid.UpdatedAt)
	return err
}

func (s *BidService) GetUserBids(username string, limit, offset int) ([]models.Bid, error) {
	query := "SELECT id, name, description, tender_id, author_type, author_id, status, version, created_at, updated_at FROM bids WHERE author_id = $1 ORDER BY name LIMIT $2 OFFSET $3"
	rows, err := s.db.Query(query, username, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Status, &bid.Version, &bid.CreatedAt, &bid.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func (s *BidService) GetBidsForTender(tenderId string, limit, offset int) ([]models.Bid, error) {
	query := "SELECT id, name, description, tender_id, author_type, author_id, status, version, created_at, updated_at FROM bids WHERE tender_id = $1 ORDER BY name LIMIT $2 OFFSET $3"
	rows, err := s.db.Query(query, tenderId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Status, &bid.Version, &bid.CreatedAt, &bid.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func (s *BidService) GetBidStatus(bidId string) (string, error) {
	var status string
	err := s.db.QueryRow("SELECT status FROM bids WHERE id = $1", bidId).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (s *BidService) UpdateBidStatus(bidId, status string) error {
	_, err := s.db.Exec("UPDATE bids SET status = $1 WHERE id = $2", status, bidId)
	return err
}

func (s *BidService) EditBid(bid *models.Bid) error {
	_, err := s.db.Exec("UPDATE bids SET name = $1, description = $2, status = $3, version = $4, updated_at = $5 WHERE id = $6",
		bid.Name, bid.Description, bid.Status, bid.Version+1, time.Now(), bid.ID)
	return err
}

func (s *BidService) SubmitBidDecision(bidId, decision string) error {
	_, err := s.db.Exec("UPDATE bids SET decision = $1 WHERE id = $2", decision, bidId)
	return err
}

func (s *BidService) SubmitBidFeedback(bidId, feedback string) error {
	_, err := s.db.Exec("UPDATE bids SET feedback = $1 WHERE id = $2", feedback, bidId)
	return err
}

func (s *BidService) RollbackBid(bidId string, version int) (*models.Bid, error) {
	var bid models.Bid
	err := s.db.QueryRow("SELECT id, name, description, tender_id, author_type, author_id, status, version, created_at, updated_at FROM bids WHERE id = $1 AND version = $2", bidId, version).
		Scan(&bid.ID, &bid.Name, &bid.Description, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Status, &bid.Version, &bid.CreatedAt, &bid.UpdatedAt)
	if err != nil {
		return nil, err
	}

	bid.Version++
	bid.UpdatedAt = time.Now()

	_, err = s.db.Exec("UPDATE bids SET name = $1, description = $2, status = $3, version = $4, updated_at = $5 WHERE id = $6",
		bid.Name, bid.Description, bid.Status, bid.Version, bid.UpdatedAt, bid.ID)
	if err != nil {
		return nil, err
	}

	return &bid, nil
}

func (s *BidService) GetBidReviews(tenderId, authorUsername string, limit, offset int) ([]models.BidReview, error) {
	query := "SELECT id, description, created_at FROM bid_reviews WHERE tender_id = $1 AND author_username = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4"
	rows, err := s.db.Query(query, tenderId, authorUsername, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.BidReview
	for rows.Next() {
		var review models.BidReview
		err := rows.Scan(&review.ID, &review.Description, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}