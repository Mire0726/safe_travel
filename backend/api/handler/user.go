package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Mire0726/safe_travel/backend/api/infrastructure/utils"
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

	ctx = utils.SetUserID(c.Request().Context(), res.ID)
	c.SetRequest(c.Request().WithContext(ctx))

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid id"})
	}

	// ユーザー削除処理
	if err := h.authUC.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"id":      id,
	})
}
