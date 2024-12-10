package handler

import (
	"fmt"
	"net/http"

	svc "github.com/Mire0726/safe_travel/backend/api/services"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateTransport(c echo.Context) error {
	ctx := c.Request().Context()
	var req svc.TransportRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request"})
	}
	id := c.Param("id")
	eventId := c.Param("eventId")

	res, err := h.transportUC.Create(ctx, req, id, eventId)
	if err != nil {
		return fmt.Errorf("イベントの作成に失敗しました: %w", err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) ListTransport(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	eventId := c.Param("eventId")

	res, err := h.transportUC.List(ctx, id, eventId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, res)
}
