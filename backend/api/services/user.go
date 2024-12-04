package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/firebase"
)

// AuthUsecase インターフェースは、認証に関するメソッドを定義します
type AuthUsecase interface {
	SignUp(ctx context.Context, req UserRequest) (*UserResponse, error)
	SignIn(ctx context.Context, req EmailPassword) (*UserResponse, error)
}

type authUC struct {
	fa   firebase.FirebaseAuth
	data datastore.Data
}

func NewAuthUC(fa firebase.FirebaseAuth, data datastore.Data) AuthUsecase {
	return &authUC{
		fa:   fa,
		data: data,
	}
}

type EmailPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest struct {
	// Email メールアドレス
	Email string `json:"email"`
	// Name 名前
	Name string `json:"name"`
	// Password パスワード
	Password string `json:"password"`
}

type UserResponse struct {
	ID           string `json:"id"`           // 内部データベース ID
	LocalID      string `json:"localId"`      // Firebase UID
	Email        string `json:"email"`        // メールアドレス
	Name         string `json:"name"`         // 名前
	IDToken      string `json:"idToken"`      // アクセストークン
	RefreshToken string `json:"refreshToken"` // リフレッシュトークン
	ExpiresIn    string `json:"expiresIn"`    // トークンの有効期限
}

func (uc *authUC) SignUp(ctx context.Context, req UserRequest) (*UserResponse, error) {
	exist, err := uc.data.ReadWriteStore().User().Exist(ctx, qm.Where("email = ?", req.Email))
	if err != nil {
		log.Println(err, "ユーザー情報の取得に失敗しました")
		return nil, fmt.Errorf("ユーザー情報の取得に失敗しました: %w", err)
	}

	if exist {
		return nil, fmt.Errorf("ユーザー情報が既に存在します")
	}

	// Firebase でユーザー作成
	firebaseUser, err := uc.fa.SignUpWithEmailPassword(ctx, req.Email, req.Password)
	if err != nil {
		log.Println(err, "firebaseのユーザー作成に失敗しました")

		return nil, err
	}

	// データベースに保存するユーザー情報
	user := &model.User{
		ID:        uuid.New().String(),
		LocalID:   firebaseUser.LocalID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// データベースにユーザー情報を保存
	if err := uc.data.ReadWriteStore().User().Insert(ctx, user); err != nil {
		if err := uc.fa.DeleteUser(ctx, firebaseUser.LocalID); err != nil {
			log.Println(err, "ユーザー情報の保存に失敗しました")
		}
		return nil, fmt.Errorf("ユーザー情報の保存に失敗しました: %w", err)
	}

	return &UserResponse{
		ID:           user.ID,              // データベース ID
		LocalID:      firebaseUser.LocalID, // Firebase UID
		Email:        req.Email,
		Name:         req.Name,
		IDToken:      firebaseUser.IDToken,
		RefreshToken: firebaseUser.RefreshToken,
		ExpiresIn:    firebaseUser.ExpiresIn,
	}, nil
}

func (uc *authUC) SignIn(ctx context.Context, req EmailPassword) (*UserResponse, error) {
	// Firebase でサインイン
	firebaseUser, err := uc.fa.SignInWithEmailPassword(ctx, req.Email, req.Password)
	if err != nil {
		log.Println(err, "firebaseのユーザーログインに失敗しました")

		return nil, fmt.Errorf("firebaseのユーザーログインに失敗しました: %w", err)
	}

	// データベースにユーザー情報を保存
	user, err := uc.data.ReadWriteStore().User().GetByEmail(ctx, req.Email)
	if err != nil {
		log.Println(err, "ユーザー情報の取得に失敗しました")
		return nil, fmt.Errorf("ユーザー情報の取得に失敗しました: %w", err)
	}

	return &UserResponse{
		ID:           user.ID,              // データベース ID
		LocalID:      firebaseUser.LocalID, // Firebase UID
		Email:        req.Email,
		Name:         user.Name,
		IDToken:      firebaseUser.IDToken,
		RefreshToken: firebaseUser.RefreshToken,
		ExpiresIn:    firebaseUser.ExpiresIn,
	}, nil
}
