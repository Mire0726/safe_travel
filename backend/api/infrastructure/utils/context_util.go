package utils

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "userID"

// SetUserID は context にユーザー ID を保存します。
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID は context からユーザー ID を取得します。
func GetUserID(ctx context.Context) (string, error) {
	value := ctx.Value(userIDKey)
	if value == nil {
		return "", errors.New("user ID not found in context")
	}
	userID, ok := value.(string)
	if !ok {
		return "", errors.New("user ID in context has invalid type")
	}
	return userID, nil
}
