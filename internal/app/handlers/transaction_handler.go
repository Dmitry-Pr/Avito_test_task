// Package handlers Description: Описывается обработчик для транзакций.
package handlers

import (
	"encoding/json"
	"merch-shop/internal/app/middleware"
	"merch-shop/internal/app/services"
	"merch-shop/internal/app/validation"
	"merch-shop/internal/pkg/errors"
	"net/http"
)

// SendCoinRequest структура запроса на отправку монет.
type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"required,min=0"`
}

var validator *validation.Validator

func init() {
	validator = validation.NewValidator()
}

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

// SendCoin отправляет монеты от одного пользователя другому.
func (h *TransactionHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	fromUserID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		jsonErr := errors.NewErrorResponse("Не авторизован")
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}

	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := errors.NewErrorResponse("Некорректный запрос")
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}

	if err := validator.ValidateSendCoinRequest(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := errors.NewErrorResponse(err.Error())
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}

	err := h.service.SendCoins(fromUserID, req.ToUser, req.Amount)
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

// GetInfo возвращает информацию о транзакциях пользователя.
func (h *TransactionHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
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

	info, err := h.service.GetUserTransactionInfo(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := errors.NewErrorResponse(err.Error())
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonErr := errors.NewErrorResponse(err.Error())
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}
}
