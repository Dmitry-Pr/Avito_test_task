// Package handlers Description: Описывает обработчики для мерча.
package handlers

import (
	"encoding/json"
	"merch-shop/internal/app/middleware"
	"merch-shop/internal/app/services"
	"merch-shop/internal/pkg/errors"
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
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		jsonErr := errors.NewErrorResponse("Не авторизован")
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}

	// Получаем item из URL-параметра
	item := ""
	pathParts := strings.Split(r.URL.Path, "/")
	item = pathParts[len(pathParts)-1]

	if item == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := errors.NewErrorResponse("Не указан предмет для покупки")
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}

	err := h.service.BuyMerch(userID, item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := errors.NewErrorResponse(err.Error())
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
