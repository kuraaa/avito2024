// cmd/main.go
package main

import (
	"avito-tender-service/internal/config"
	"avito-tender-service/internal/db"
	"avito-tender-service/internal/handlers"
	"avito-tender-service/internal/services"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := db.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	tenderService := services.NewTenderService(db)
	tenderHandler := handlers.NewTenderHandler(tenderService)

	bidService := services.NewBidService(db)
	bidHandler := handlers.NewBidHandler(bidService)

	r := chi.NewRouter()

	r.Get("/api/ping", handlers.PingHandler)

	r.Get("/api/tenders", tenderHandler.GetTenders)
	r.Post("/api/tenders/new", tenderHandler.CreateTender)
	r.Get("/api/tenders/my", tenderHandler.GetUserTenders)
	r.Get("/api/tenders/{tenderId}/status", tenderHandler.GetTenderStatus)
	r.Put("/api/tenders/{tenderId}/status", tenderHandler.UpdateTenderStatus)
	r.Patch("/api/tenders/{tenderId}/edit", tenderHandler.EditTender)
	r.Put("/api/tenders/{tenderId}/rollback/{version}", tenderHandler.RollbackTender)

	r.Post("/api/bids/new", bidHandler.CreateBid)
	r.Get("/api/bids/my", bidHandler.GetUserBids)
	r.Get("/api/bids/{tenderId}/list", bidHandler.GetBidsForTender)
	r.Get("/api/bids/{bidId}/status", bidHandler.GetBidStatus)
	r.Put("/api/bids/{bidId}/status", bidHandler.UpdateBidStatus)
	r.Patch("/api/bids/{bidId}/edit", bidHandler.EditBid)
	r.Put("/api/bids/{bidId}/submit_decision", bidHandler.SubmitBidDecision)
	r.Put("/api/bids/{bidId}/feedback", bidHandler.SubmitBidFeedback)
	r.Put("/api/bids/{bidId}/rollback/{version}", bidHandler.RollbackBid)
	r.Get("/api/bids/{tenderId}/reviews", bidHandler.GetBidReviews)

	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
