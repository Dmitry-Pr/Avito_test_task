// Package di Description: этот пакет содержит зависимости приложения.
package di

import (
	"merch-shop/internal/app/handlers"
	"merch-shop/internal/app/repositories"
	"merch-shop/internal/app/services"

	"gorm.io/gorm"
)

// Dependencies содержит зависимости приложения.
type Dependencies struct {
	DB                 *gorm.DB
	MerchRepo          repositories.MerchRepositoryInterface
	MerchService       services.MerchServiceInterface
	MerchHandler       handlers.MerchHandlerInterface
	UserRepo           repositories.UserRepositoryInterface
	UserService        services.UserServiceInterface
	UserHandler        handlers.UserHandlerInterface
	InventoryRepo      repositories.InventoryRepositoryInterface
	TransactionRepo    repositories.TransactionRepositoryInterface
	TransactionService services.TransactionServiceInterface
	TransactionHandler handlers.TransactionHandlerInterface
}

// BuildDependencies создает новые зависимости приложения.
func BuildDependencies(db *gorm.DB) *Dependencies {
	userRepo := repositories.NewUserRepository(db)
	merchRepo := repositories.NewMerchRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	userService := services.NewUserService(userRepo)
	transactionService := services.NewTransactionService(userRepo, transactionRepo, merchRepo, inventoryRepo)
	merchService := services.NewMerchService(merchRepo, userRepo, transactionRepo, inventoryRepo)

	userHandler := handlers.NewUserHandler(userService)
	merchHandler := handlers.NewMerchHandler(merchService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	return &Dependencies{
		DB:                 db,
		MerchRepo:          merchRepo,
		MerchService:       merchService,
		MerchHandler:       merchHandler,
		UserRepo:           userRepo,
		UserService:        userService,
		UserHandler:        userHandler,
		TransactionRepo:    transactionRepo,
		InventoryRepo:      inventoryRepo,
		TransactionService: transactionService,
		TransactionHandler: transactionHandler,
	}
}
