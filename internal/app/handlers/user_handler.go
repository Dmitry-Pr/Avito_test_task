package handlers

import (
	"encoding/json"
	"merch-shop/internal/app/services"
	"net/http"
)

type UserHandlerInterface interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Ошибка декодирования запроса", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}
	token, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(AuthResponse{Token: token})
	if err != nil {
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
		return
	}
}
