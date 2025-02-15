package handlers

import (
	"encoding/json"
	"merch-shop/internal/app/services"
	"net/http"
)

type IUserHandler interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	userService services.IUserService
}

func NewUserHandler(userService services.IUserService) *UserHandler {
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
	json.NewDecoder(r.Body).Decode(&req)
	token, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}
