package handlers

import (
	"merch-shop/internal/app/services"
	"net/http"
	"strings"
)

type MerchHandlerInterface interface {
	BuyMerch(w http.ResponseWriter, r *http.Request)
}

type MerchHandler struct {
	service services.MerchServiceInterface
}

func NewMerchHandler(service services.MerchServiceInterface) MerchHandlerInterface {
	return &MerchHandler{service: service}
}

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
