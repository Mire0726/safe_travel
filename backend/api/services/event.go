package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// EventUsecase インターフェースは、イベントに関するメソッドを定義します
type EventUsecase interface {
	Create(ctx context.Context, req EventRequest, userID string) (*EventResponse, error)
	List(ctx context.Context, userID string) ([]*model.Event, error)
	Delete(ctx context.Context, id, eventId string) error
}

type eventUC struct {
	data datastore.Data
}

func NewEventUC(data datastore.Data) EventUsecase {
	return &eventUC{
		data: data,
	}
}

// EventRequest イベントリクエスト
type EventRequest struct {
	// Name イベント名
	Name string `json:"name"`
}

// EventResponse イベントレスポンス
type EventResponse struct {
	// ID イベント ID
	ID string `json:"id"`
	// Name イベント名
	Name string `json:"name"`
	// CreatedBy 作成者
	CreatedBy string `json:"createdBy"`
}

func (uc *eventUC) Create(ctx context.Context, req EventRequest, userID string) (*EventResponse, error) {
	// イベントを作成
	event := &model.Event{
		Name:      req.Name,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.data.ReadWriteStore().Event().Insert(ctx, event); err != nil {
		return nil, fmt.Errorf("イベントの作成に失敗しました: %w", err)
	}

	return nil, nil
}

func (uc *eventUC) List(ctx context.Context, userID string) ([]*model.Event, error) {
	events, err := uc.data.ReadWriteStore().Event().List(ctx, qm.Where("created_by = ?", userID))
	if err != nil {
		return nil, fmt.Errorf("イベントの取得に失敗しました: %w", err)
	}

	return events, nil
}

func (uc *eventUC) Delete(ctx context.Context, id, eventId string) error {

	// イベントが存在するか確認
	exist, err := uc.data.ReadWriteStore().Event().Exist(ctx, qm.Where("id = ? AND created_by = ?", eventId, id))
	if err != nil {
		return fmt.Errorf("イベントの存在確認に失敗しました: %w", err)
	}
	if !exist {
		return fmt.Errorf("イベントが存在しません")
	}

	// イベントを削除
	if err := uc.data.ReadWriteStore().Event().Delete(ctx, eventId); err != nil {
		return fmt.Errorf("イベントの削除に失敗しました: %w", err)
	}

	return nil
}
