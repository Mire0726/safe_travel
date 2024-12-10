package handler

import (
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/firebase"
	"github.com/Mire0726/safe_travel/backend/api/services"
)

type Handler struct {
	authUC  services.AuthUsecase
	eventUC services.EventUsecase
	transportUC services.TransportUsecase
}

func NewHandler(fa firebase.FirebaseAuth, data datastore.Data) *Handler {
	return &Handler{
		authUC:  services.NewAuthUC(fa, data),
		eventUC: services.NewEventUC(data),
		transportUC: services.NewTransportUC(data),
	}
}
