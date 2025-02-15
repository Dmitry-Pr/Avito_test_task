package server

import (
	"log"
	"merch-store/internal/app/di"
	"merch-store/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, container *di.Dependencies) *Server {
	mux := http.NewServeMux()

	// Используем обработчики из DI-контейнера
	mux.HandleFunc("/merch", container.MerchHandler.GetMerch)

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
