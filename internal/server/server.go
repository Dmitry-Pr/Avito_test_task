// Package server Description: Описывает сервер приложения
package server

import (
	"context"
	"errors"
	"log"
	"merch-shop/internal/app/di"
	"merch-shop/internal/app/middleware"
	"merch-shop/internal/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server описывает сервер приложения
type Server struct {
	httpServer *http.Server
}

// NewServer создает новый сервер приложения
func NewServer(cfg *config.Config, dependencies *di.Dependencies) *Server {
	mux := http.NewServeMux()

	mux.Handle("/api/auth",
		middleware.MethodMiddleware(http.HandlerFunc(dependencies.UserHandler.Authenticate),
			http.MethodPost))
	mux.Handle("/api/buy/{item}",
		middleware.MethodMiddleware(http.HandlerFunc(dependencies.MerchHandler.BuyMerch),
			http.MethodPost))
	mux.Handle("/api/info",
		middleware.MethodMiddleware(http.HandlerFunc(dependencies.TransactionHandler.GetInfo),
			http.MethodGet))
	mux.Handle("/api/sendCoin",
		middleware.MethodMiddleware(http.HandlerFunc(dependencies.TransactionHandler.SendCoin),
			http.MethodPost))

	handlerWithMiddleware := middleware.AuthMiddleware(mux)
	handlerWithMiddleware = middleware.LogsMiddleware(handlerWithMiddleware)

	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.Server.Address,
			Handler:           handlerWithMiddleware,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

// Run запускает сервер
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
		log.Printf("Ошибка при завершении сервера: %v", err)
	}

	log.Println("Сервер корректно завершил работу")
	return nil
}

// HTTPHandler - метод получения хендлера для тестов
func (s *Server) HTTPHandler() http.Handler {
	return s.httpServer.Handler
}
