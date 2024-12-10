package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// TransportUsecase インターフェースは、トランスポートに関するメソッドを定義します
type TransportUsecase interface {
	Create(ctx context.Context, req TransportRequest, userID, eventID string) (*TransportResponse, error)
	List(ctx context.Context, userID, eventID string) ([]*model.Transport, error)
}

type transportUC struct {
	data datastore.Data
}

func NewTransportUC(data datastore.Data) TransportUsecase {
	return &transportUC{data: data}
}

// TransportRequest 交通機関リクエスト
type TransportRequest struct {
	TransportType  string    `json:"transportType"`
	Memo           string    `json:"memo"`
	StartLocation  string    `json:"startLocation"`
	ArriveLocation string    `json:"arriveLocation"`
	StartAt        time.Time `json:"startAt"`
	ArriveAt       time.Time `json:"arriveAt"`
}

// TransportResponse 交通機関レスポンス
type TransportResponse struct {
	ID             string                        `json:"id"`
	TransportType  model.TransportsTransportType `json:"transportType"`
	Memo           string                        `json:"memo"`
	StartLocation  string                        `json:"startLocation"`
	ArriveLocation string                        `json:"arriveLocation"`
	StartAt        time.Time                     `json:"startAt"`
	ArriveAt       time.Time                     `json:"arriveAt"`
}

func (uc *transportUC) Create(ctx context.Context, req TransportRequest, userID, eventID string) (*TransportResponse, error) {
	exist, err := uc.data.ReadWriteStore().Event().Exist(ctx, qm.Where("id = ? AND created_by = ?", eventID, userID))
	if err != nil {
		return nil, fmt.Errorf("イベントの存在確認に失敗しました: %w", err)
	}

	if !exist {
		return nil, fmt.Errorf("イベントが見つかりません")
	}

	// 交通機関を作成
	transport := &model.Transport{
		ID:             uuid.New().String(),
		TransportType:  model.TransportsTransportType(req.TransportType),
		Memo:           req.Memo,
		StartLocation:  req.StartLocation,
		ArriveLocation: req.ArriveLocation,
		StartAt:        req.StartAt,
		ArriveAt:       req.ArriveAt,
		EventID:        eventID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.data.ReadWriteStore().Transport().Insert(ctx, transport); err != nil {
		return nil, fmt.Errorf("交通機関の作成に失敗しました: %w", err)
	}

	return &TransportResponse{
		ID:             transport.ID,
		TransportType:  transport.TransportType,
		Memo:           transport.Memo,
		StartLocation:  transport.StartLocation,
		ArriveLocation: transport.ArriveLocation,
		StartAt:        transport.StartAt,
		ArriveAt:       transport.ArriveAt,
	}, nil
}

func (uc *transportUC) List(ctx context.Context, userID, eventID string) ([]*model.Transport, error) {
	exist, err := uc.data.ReadWriteStore().Event().Exist(ctx, qm.Where("id = ? AND created_by = ?", eventID, userID))
	if err != nil {
		return nil, fmt.Errorf("イベントの存在確認に失敗しました: %w", err)
	}

	if !exist {
		return nil, fmt.Errorf("イベントが見つかりません")
	}

	transports, err := uc.data.ReadWriteStore().Transport().List(ctx, qm.Where("event_id = ?", eventID))
	if err != nil {
		return nil, fmt.Errorf("交通機関の取得に失敗しました: %w", err)
	}

	return transports, nil
}

func (uc *transportUC) Delete(ctx context.Context, userID, eventID, transportID string) error {
	// 交通機関が存在するか確認
	exist, err := uc.data.ReadWriteStore().Transport().Exist(ctx, qm.Where("id = ? AND event_id = ?", transportID, eventID))
	if err != nil {
		return fmt.Errorf("交通機関の存在確認に失敗しました: %w", err)
	}
	if !exist {
		return fmt.Errorf("交通機関が存在しません")
	}

	if err := uc.data.ReadWriteStore().Transport().Delete(ctx, transportID); err != nil {
		return fmt.Errorf("交通機関の削除に失敗しました: %w", err)
	}

	return nil
}
