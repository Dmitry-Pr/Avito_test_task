package server

import (
	"context"
	"errors"
	"log"
	"merch-shop/internal/app/di"
	"merch-shop/internal/app/middleware"
	"merch-shop/internal/config"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, dependencies *di.Dependencies) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth", dependencies.UserHandler.Authenticate)
	mux.HandleFunc("/api/merch", dependencies.MerchHandler.GetMerch)

	handlerWithMiddleware := middleware.AuthMiddleware(mux)
	handlerWithMiddleware = middleware.LogsMiddleware(mux)

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Server.Address,
			Handler: handlerWithMiddleware,
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
