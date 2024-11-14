package services

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/firebase"
	"gorm.io/gorm/logger"
)

// AuthUsecase インターフェースは、認証に関するメソッドを定義します
type AuthUsecase interface {
	SignUp(ctx context.Context, email, name, password string) (*auth.UserRecord, error)
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

type SignRequest struct {
	// Email メールアドレス
	Email string `json:"email"`
	// Name 名前
	Name string `json:"name"`
	// Password パスワード
	Password string `json:"password"`
}

func (uc *authUC) SignUp(ctx context.Context, email, name, password string) (*auth.UserRecord, error) {
	if err := uc.data.ReadWriteStore().User().Insert(ctx, &model.User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}); err != nil {
		logger.Default.LogMode(logger.Error)
		fmt.Println(err, "ユーザー情報のDB保存に失敗しました")
		return nil, err
	}

	res, err := uc.fa.CreateUser(ctx, email, password)
	if err != nil {
		logger.Default.LogMode(logger.Error)
		fmt.Println(err, "firebaseのユーザー作成に失敗しました")

		return nil, err
	}
	return res, nil
}
