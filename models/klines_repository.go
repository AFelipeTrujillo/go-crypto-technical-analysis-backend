package models

import (
	"gorm.io/gorm"
)

type KlinesRepositoryInterface interface {
	GetAll(symbol string, startTimestamp, endTimestamp int64) ([]Kline, int, error)
}

type KlinesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) KlinesRepositoryInterface {
	return &KlinesRepository{
		db: db,
	}
}

func (r *KlinesRepository) GetAll(symbol string, startTimestamp, endTimestamp int64) ([]Kline, int, error) {
	var klines []Kline

	query := r.db.Model(&Kline{}).
		Joins("JOIN symbols ON symbols.id = klines_1h.symbol_id").
		Where("symbols.name = ?", symbol)

	if startTimestamp > 0 {
		query = query.Where("klines_1h.open_time >= ?", startTimestamp)
	}

	if endTimestamp > 0 {
		query = query.Where("klines_1h.open_time <= ?", endTimestamp)
	}

	if err := query.Order("klines_1h.open_time DESC").Find(&klines).Error; err != nil {
		return nil, 0, err
	}
	return klines, len(klines), nil
}
