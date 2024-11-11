package services

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/firebase"
	"gorm.io/gorm/logger"
)

// AuthUsecase インターフェースは、認証に関するメソッドを定義します
type AuthUsecase interface {
	SignUp(ctx context.Context, email, password string) (*auth.UserRecord, error)
}

type authUC struct {
	fa *firebase.FirebaseAuth
}

func NewAuthUC(fa *firebase.FirebaseAuth) AuthUsecase {
	return &authUC{fa: fa}
}

type SignRequest struct {
	// Email メールアドレス
	Email string `json:"email"`

	// Password パスワード
	Password string `json:"password"`
}

func (uc *authUC) SignUp(ctx context.Context, email, password string) (*auth.UserRecord, error) {
	res, err := uc.fa.CreateUser(ctx, email, password)
	if err != nil {
		logger.Default.LogMode(logger.Error)

		return nil, err
	}

	return res, nil
}
