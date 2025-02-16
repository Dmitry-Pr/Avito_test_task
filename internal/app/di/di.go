package di

import (
	"merch-shop/internal/app/handlers"
	"merch-shop/internal/app/repositories"
	"merch-shop/internal/app/services"

	"gorm.io/gorm"
)

type Dependencies struct {
	DB              *gorm.DB
	MerchRepo       repositories.MerchRepositoryInterface
	MerchService    services.MerchServiceInterface
	MerchHandler    handlers.MerchHandlerInterface
	UserRepo        repositories.UserRepositoryInterface
	UserService     services.UserServiceInterface
	UserHandler     handlers.UserHandlerInterface
	TransactionRepo repositories.TransactionRepositoryInterface
}

func BuildDependencies(db *gorm.DB) *Dependencies {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	transactionRepo := repositories.NewTransactionRepository(db)

	merchRepo := repositories.NewMerchRepository(db)
	merchService := services.NewMerchService(merchRepo, userRepo, transactionRepo)
	merchHandler := handlers.NewMerchHandler(merchService)

	return &Dependencies{
		DB:              db,
		MerchRepo:       merchRepo,
		MerchService:    merchService,
		MerchHandler:    merchHandler,
		UserRepo:        userRepo,
		UserService:     userService,
		UserHandler:     userHandler,
		TransactionRepo: transactionRepo,
	}
}
