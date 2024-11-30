package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/firebase"
)

// AuthUsecase インターフェースは、認証に関するメソッドを定義します
type AuthUsecase interface {
	SignUp(ctx context.Context, email, name, password string) (*SignUpResponse, error)
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

type SignUpRequest struct {
	// Email メールアドレス
	Email string `json:"email"`
	// Name 名前
	Name string `json:"name"`
	// Password パスワード
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID           string `json:"id"`           // 内部データベース ID
	LocalID      string `json:"localId"`      // Firebase UID
	Email        string `json:"email"`        // メールアドレス
	Name         string `json:"name"`         // 名前
	IDToken      string `json:"idToken"`      // アクセストークン
	RefreshToken string `json:"refreshToken"` // リフレッシュトークン
	ExpiresIn    string `json:"expiresIn"`    // トークンの有効期限
}

func (uc *authUC) SignUp(ctx context.Context, email, name, password string) (*SignUpResponse, error) {
	// Firebase でユーザー作成
	firebaseUser, err := uc.fa.SignUpWithEmailPassword(ctx, email, password)
	if err != nil {
		log.Println(err, "firebaseのユーザー作成に失敗しました")

		return nil, err
	}

	// データベースに保存するユーザー情報
	user := &model.User{
		ID:        uuid.New().String(),
		LocalID:   firebaseUser.LocalID,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// データベースにユーザー情報を保存
	if err := uc.data.ReadWriteStore().User().Insert(ctx, user); err != nil {
		uc.fa.DeleteUser(ctx, firebaseUser.LocalID)
		return nil, fmt.Errorf("ユーザー情報の保存に失敗しました: %w", err)
	}

	return &SignUpResponse{
		ID:           user.ID,              // データベース ID
		LocalID:      firebaseUser.LocalID, // Firebase UID
		Email:        email,
		Name:         name,
		IDToken:      firebaseUser.IDToken,
		RefreshToken: firebaseUser.RefreshToken,
		ExpiresIn:    firebaseUser.ExpiresIn,
	}, nil
}

// func (uc *authUC) SignUp(ctx context.Context, email, name, password string) (*SignUpResponse, error) {
// 	user := &model.User{
// 		ID:    uuid.New().String(),
// 		Name:  name,
// 		Email: email,
// 	}

// 	if err := uc.data.ReadWriteStore().User().Insert(ctx, user); err != nil {
// 		fmt.Println(err, "ユーザー情報の保存に失敗しました")

// 		return nil, err
// 	}

// 	res, err := uc.fa.CreateUser(ctx, email, password)
// 	if err != nil {
// 		fmt.Println(err, "firebaseのユーザー作成に失敗しました")

// 		return nil, err
// 	}
// 	return &SignUpResponse{
// 		Res:  res,
// 		User: user,
// 	}, nil
// }
