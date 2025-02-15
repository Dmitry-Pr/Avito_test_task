package di

import (
	"gorm.io/gorm"
	"merch-shop/internal/app/handlers"
	"merch-shop/internal/app/repositories"
	"merch-shop/internal/app/services"
)

type Dependencies struct {
	DB           *gorm.DB
	MerchRepo    repositories.MerchRepositoryInterface
	MerchService services.MerchServiceInterface
	MerchHandler handlers.MerchHandlerInterface
	UserRepo     repositories.IUserRepository
	UserService  services.IUserService
	UserHandler  handlers.IUserHandler
}

func BuildDependencies(db *gorm.DB) *Dependencies {
	merchRepo := repositories.NewMerchRepository(db)
	merchService := services.NewMerchService(merchRepo)
	merchHandler := handlers.NewMerchHandler(merchService)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	return &Dependencies{
		DB:           db,
		MerchRepo:    merchRepo,
		MerchService: merchService,
		MerchHandler: merchHandler,
		UserRepo:     userRepo,
		UserService:  userService,
		UserHandler:  userHandler,
	}
}
