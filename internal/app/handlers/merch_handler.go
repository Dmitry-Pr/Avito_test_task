package handlers

import (
	"encoding/json"
	"merch-shop/internal/app/services"
	"net/http"
)

type MerchHandlerInterface interface {
	GetMerch(w http.ResponseWriter, r *http.Request)
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
