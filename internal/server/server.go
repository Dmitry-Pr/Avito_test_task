package server

import (
	"context"
	"errors"
	"log"
	"merch-store/internal/app/di"
	"merch-store/internal/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, container *di.Dependencies) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/merch", container.MerchHandler.GetMerch)

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Server.Address,
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	log.Println("Запускаем сервер на", s.httpServer.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("Сервер остановлен")
			} else {
				log.Fatalf("Ошибка запуска сервера: %v", err)
			}
		}
	}()

	<-quit
	log.Println("Получен сигнал завершения, останавливаем сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}

	log.Println("Сервер корректно завершил работу")
	return nil
}
