package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware は認証ミドルウェアです。
func AuthMiddleware(client *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// "Bearer "を除去してトークンを取得
			idToken := strings.TrimPrefix(authHeader, "Bearer ")
			// Firebaseでトークンを検証してユーザーIDを取得
			token, err := client.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			// TODO: コンテキストに入れるUserIDはFirebaseのUIDではなく、データベースのIDにする =>カスタムクレイム？
			ctx := SetUserID(c.Request().Context(), token.UID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

type contextKey string

const userIDKey contextKey = "userID"

// SetUserID は context にユーザー ID を保存します。
func SetUserID(ctx context.Context, userID string) context.Context {
	log.Println("SetUserID", userID)
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID は context からユーザー ID を取得します。
func GetUserID(ctx context.Context) (string, error) {
	value := ctx.Value(userIDKey)
	if value == nil {
		return "", fmt.Errorf("user ID not found in context")
	}
	userID, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("failed to assert user ID")
	}
	return userID, nil
}
