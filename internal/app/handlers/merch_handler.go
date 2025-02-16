package handlers

import (
	"encoding/json"
	"merch-shop/internal/app/services"
	"net/http"
	"strings"
)

type MerchHandlerInterface interface {
	GetMerch(w http.ResponseWriter, r *http.Request)
	BuyMerch(w http.ResponseWriter, r *http.Request)
}

type MerchHandler struct {
	service services.MerchServiceInterface
}

func NewMerchHandler(service services.MerchServiceInterface) MerchHandlerInterface {
	return &MerchHandler{service: service}
}

func (h *MerchHandler) GetMerch(w http.ResponseWriter, r *http.Request) {
	merch, err := h.service.GetAllMerch()
	if err != nil {
		http.Error(w, "Ошибка получения товаров", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(merch); err != nil {
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
	}
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
