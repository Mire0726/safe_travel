package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

// UserContextKey はコンテキストにユーザー情報を格納するためのキー型です
type UserContextKey string

const userContextKey UserContextKey = "user"

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
func (fa *FirebaseAuth) VerifyToken(ctx context.Context, authHeader string) (context.Context, error) {
	if authHeader == "" {
		return nil, echo.NewHTTPError(401, "missing authorization header")
	}

	// "Bearer "を除去してトークンを取得
	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := fa.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, echo.NewHTTPError(401, "invalid token")
	}

	// コンテキストにトークン情報をセット
	ctx = context.WithValue(ctx, userContextKey, token)

	return ctx, nil
}

// GetUserID はコンテキストからユーザーIDを取得するヘルパー関数です
func GetUserID(c echo.Context) string {
	token := c.Get("user").(*auth.Token)
	return token.UID
}

// GetUser はユーザー情報を取得します
func (fa *FirebaseAuth) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := fa.auth.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (fa *FirebaseAuth) DeleteUser(ctx context.Context, uid string) error {
	if err := fa.auth.DeleteUser(ctx, uid); err != nil {
		return err
	}

	return nil
}

type SignUpResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

type signUpRequestWithEmailPassword struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

func (fa *FirebaseAuth) SignUpWithEmailPassword(ctx context.Context, email, password string) (*SignUpResponse, error) {
	firebaseAPIKey := os.Getenv("FIREBASE_API_KEY")

	reqBody := &signUpRequestWithEmailPassword{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=%s", firebaseAPIKey)

	signUpResponse := &SignUpResponse{}
	if err := fa.callPost(ctx, url, reqBody, &signUpResponse); err != nil {
		fmt.Println(err, "infra:firebaseのユーザー作成に失敗しました")
		
		return nil, err
	}

	return signUpResponse, nil
}

func (fa *FirebaseAuth) callPost(ctx context.Context, url string, reqBody any, respBody interface{}) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println(err, "infra:firebaseのリクエストボディのJSON変換に失敗しました")

		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err, "infra:firebaseのリクエスト作成に失敗しました")

		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err, "infra:firebaseのリクエスト送信に失敗しました")

		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err, "infra:firebaseのレスポンスボディの読み込みに失敗しました")

		return err
	}

	if err := json.Unmarshal(body, &respBody); err != nil {
		fmt.Println(err, "infra:firebaseのレスポンスボディのJSON変換に失敗しました")

		return err
	}

	return nil
}
