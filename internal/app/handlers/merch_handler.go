// Package handlers Description: Описывает обработчики для мерча.
package handlers

import (
	"merch-shop/internal/app/services"
	"net/http"
	"strings"
)

// MerchHandlerInterface описывает обработчик для мерча.
type MerchHandlerInterface interface {
	BuyMerch(w http.ResponseWriter, r *http.Request)
}

// MerchHandler обработчик для мерча.
type MerchHandler struct {
	service services.MerchServiceInterface
}

// NewMerchHandler создает новый обработчик для мерча.
func NewMerchHandler(service services.MerchServiceInterface) MerchHandlerInterface {
	return &MerchHandler{service: service}
}

// BuyMerch купить мерч.
func (h *MerchHandler) BuyMerch(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	// Получаем item из URL-параметра
	item := ""
	pathParts := strings.Split(r.URL.Path, "/")
	item = pathParts[len(pathParts)-1]

	if item == "" {
		http.Error(w, "Не указан предмет для покупки", http.StatusBadRequest)
		return
	}

	err := h.service.BuyMerch(userID, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
