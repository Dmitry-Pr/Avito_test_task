// Package handlers Description: Описывает обработчики запросов для пользователей.
package handlers

import (
	"encoding/json"
	"merch-shop/internal/app/services"
	"merch-shop/internal/pkg/errors"
	"net/http"
)

// UserHandlerInterface описывает обработчик запросов для пользователей.
type UserHandlerInterface interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

// UserHandler обработчик запросов для пользователей.
type UserHandler struct {
	userService services.UserServiceInterface
}

// NewUserHandler создает новый обработчик запросов для пользователей.
func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

// AuthRequest структура запроса на аутентификацию.
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse структура ответа на аутентификацию.
type AuthResponse struct {
	Token string `json:"token"`
}

// Authenticate аутентификация пользователя.
func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := errors.NewErrorResponse("Ошибка декодирования запроса")
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}
	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := errors.NewErrorResponse("Некорректные данные")
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}
	token, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		jsonErr := errors.NewErrorResponse("Не авторизован")
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(AuthResponse{Token: token})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonErr := errors.NewErrorResponse("Ошибка кодирования ответа")
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
		return
	}
}
