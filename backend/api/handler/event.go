package handler

import (
	"net/http"

	svc "github.com/Mire0726/safe_travel/backend/api/services"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateEvent(c echo.Context) error {
	ctx := c.Request().Context()
	var req svc.EventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request"})
	}
	id := c.Param("id")

	res, err := h.eventUC.Create(ctx, req, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) ListEvent(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	res, err := h.eventUC.List(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, res)
}
