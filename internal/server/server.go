package server

import (
	"database/sql"
	"log"
	"merch-store/internal/app/handlers"
	"merch-store/internal/app/repositories"
	"merch-store/internal/app/services"
	"merch-store/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	mux := http.NewServeMux()

	// Создаем репозиторий, сервис и хендлер
	merchRepo := repositories.NewMerchRepository(db)
	merchService := services.NewMerchService(merchRepo)
	merchHandler := handlers.NewMerchHandler(merchService)

	mux.HandleFunc("/merch", merchHandler.GetMerch)

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Server.Address,
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	log.Println("Сервер запущен на", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
