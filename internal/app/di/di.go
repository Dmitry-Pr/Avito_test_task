package di

import (
	"database/sql"
	"merch-store/internal/app/handlers"
	"merch-store/internal/app/repositories"
	"merch-store/internal/app/services"
)

type Dependencies struct {
	DB           *sql.DB
	MerchRepo    *repositories.MerchRepository
	MerchService *services.MerchService
	MerchHandler *handlers.MerchHandler
}

func BuildDependencies(db *sql.DB) *Dependencies {
	merchRepo := repositories.NewMerchRepository(db)
	merchService := services.NewMerchService(merchRepo)
	merchHandler := handlers.NewMerchHandler(merchService)

	return &Dependencies{
		DB:           db,
		MerchRepo:    merchRepo,
		MerchService: merchService,
		MerchHandler: merchHandler,
	}
}
