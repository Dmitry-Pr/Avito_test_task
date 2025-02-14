package handlers

import (
	"encoding/json"
	"merch-store/internal/app/services"
	"net/http"
)

type MerchHandler struct {
	service *services.MerchService
}

func NewMerchHandler(service *services.MerchService) *MerchHandler {
	return &MerchHandler{service: service}
}

func (h *MerchHandler) GetMerch(w http.ResponseWriter, r *http.Request) {
	merch, err := h.service.GetAllMerch()
	if err != nil {
		http.Error(w, "Ошибка получения товаров", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(merch)
	if err != nil {
		return
	}
}
