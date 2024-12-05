package firebase

import (

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	fa *FirebaseAuth
}

func NewAuthMiddleware(fa *FirebaseAuth) *AuthMiddleware {
	return &AuthMiddleware{fa: fa}
}

func (m *AuthMiddleware) VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := m.fa.VerifyToken(
			c.Request().Context(),
			c.Request().Header.Get("Authorization"),
		)
		if err != nil {
			return err
		}

		c.Set("user", token)
		return next(c)
	}
}