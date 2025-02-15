package di

import (
	"database/sql"
	"merch-shop/internal/app/handlers"
	"merch-shop/internal/app/repositories"
	"merch-shop/internal/app/services"
)

type Dependencies struct {
	DB           *sql.DB
	MerchRepo    repositories.MerchRepositoryInterface
	MerchService services.MerchServiceInterface
	MerchHandler handlers.MerchHandlerInterface
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
