package user

import (
	"context"
	"fmt"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
)

func (u *user) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := model.Users(
		model.UserWhere.Email.EQ(email),
	).One(ctx, u.dbClient)
	if err != nil {
		return nil, fmt.Errorf("error executing user.GetByEmail: %w", err)
	}

	return user, nil
}
