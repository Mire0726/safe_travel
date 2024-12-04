package handler

import (
	"net/http"

	svc "github.com/Mire0726/safe_travel/backend/api/services"
	"github.com/labstack/echo/v4"
)

func (h *Handler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	var req svc.UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request"})
	}

	res, err := h.authUC.SignUp(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) SignIn(c echo.Context) error {
	ctx := c.Request().Context()
	var req svc.EmailPassword
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request"})
	}

	res, err := h.authUC.SignIn(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, res)
}
