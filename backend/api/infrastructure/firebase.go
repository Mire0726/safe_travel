package infrastructure

import (
	"context"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	app  *firebase.App
	auth *auth.Client
}

// NewFirebaseAuth はFirebase認証クライアントを初期化します
func NewFirebaseAuth() (*FirebaseAuth, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./api/config/firebase-key.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseAuth{
		app:  app,
		auth: auth,
	}, nil
}

// VerifyToken はJWTトークンを検証します
func (fa *FirebaseAuth) VerifyToken(ctx context.Context, authHeader string) (*auth.Token, error) {
	if authHeader == "" {
		return nil, echo.NewHTTPError(401, "missing authorization header")
	}

	// "Bearer "を除去してトークンを取得
	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := fa.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, echo.NewHTTPError(401, "invalid token")
	}

	return token, nil
}

// GetUserID はコンテキストからユーザーIDを取得するヘルパー関数です
func GetUserID(c echo.Context) string {
	token := c.Get("user").(*auth.Token)
	return token.UID
}
