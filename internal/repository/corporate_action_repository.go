package repository

import (
	"context"
	"time"

	"echobackend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CorporateActionRepository persists the IDX dividend & RUPS calendar in Postgres.
type CorporateActionRepository interface {
	// ExistsInRange reports whether any corporate action rows already fall
	// within [from, to]. Used to decide whether an IDX fetch is needed.
	ExistsInRange(ctx context.Context, from, to time.Time) (bool, error)
	FindByDateRange(ctx context.Context, from, to time.Time) ([]model.CorporateAction, error)
	// UpsertMany inserts actions, updating mutable fields on conflict of
	// (symbol, type, event_date).
	UpsertMany(ctx context.Context, actions []model.CorporateAction) error
}

type corporateActionRepository struct {
	db *gorm.DB
}

func NewCorporateActionRepository(db *gorm.DB) CorporateActionRepository {
	return &corporateActionRepository{db: db}
}

func (r *corporateActionRepository) ExistsInRange(ctx context.Context, from, to time.Time) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.CorporateAction{}).
		Where("event_date BETWEEN ? AND ?", from, to).
		Count(&count).Error
	return count > 0, err
}

func (r *corporateActionRepository) FindByDateRange(ctx context.Context, from, to time.Time) ([]model.CorporateAction, error) {
	var actions []model.CorporateAction
	err := r.db.WithContext(ctx).
		Where("event_date BETWEEN ? AND ?", from, to).
		Order("event_date ASC").
		Find(&actions).Error
	return actions, err
}

func (r *corporateActionRepository) UpsertMany(ctx context.Context, actions []model.CorporateAction) error {
	if len(actions) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "symbol"}, {Name: "type"}, {Name: "event_date"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "pay_date", "amount", "currency", "note", "market", "updated_at"}),
	}).Create(&actions).Error
}
