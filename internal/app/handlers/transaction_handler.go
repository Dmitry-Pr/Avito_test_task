// Package handlers Description: Описывается обработчик для транзакций.
package handlers

import (
	"encoding/json"
	"merch-shop/internal/app/services"
	"merch-shop/internal/app/validation"
	"net/http"
)

// TransactionHandlerInterface описывает обработчик для транзакций.
type TransactionHandlerInterface interface {
	SendCoin(w http.ResponseWriter, r *http.Request)
	GetInfo(w http.ResponseWriter, r *http.Request)
}

// TransactionHandler обработчик для транзакций.
type TransactionHandler struct {
	service services.TransactionServiceInterface
}

// NewTransactionHandler создает новый обработчик для транзакций.
func NewTransactionHandler(service services.TransactionServiceInterface) TransactionHandlerInterface {
	return &TransactionHandler{service: service}
}

// SendCoinRequest структура запроса на отправку монет.
type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"required,min=0"`
}

var validator *validation.Validator // Создаем экземпляр валидатора

func init() {
	validator = validation.NewValidator()
}

// SendCoin отправляет монеты от одного пользователя другому.
func (h *TransactionHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	fromUserID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateSendCoinRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.SendCoins(fromUserID, req.ToUser, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetInfo возвращает информацию о транзакциях пользователя.
func (h *TransactionHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	info, err := h.service.GetUserTransactionInfo(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
