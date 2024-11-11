package handler

import "github.com/Mire0726/safe_travel/backend/api/services"

type Handler struct {
	authUC services.AuthUsecase
}

func NewHandler(authUC services.AuthUsecase) *Handler {
	return &Handler{authUC: authUC}
}