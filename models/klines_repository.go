package models

import (
	"fmt"

	"gorm.io/gorm"
)

type KlinesRepositoryInterface interface {
	GetAll(startTimestamp, endTimestamp int64) ([]Kline, int, error)
}

type KlinesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) KlinesRepositoryInterface {
	return &KlinesRepository{
		db: db,
	}
}

func (r *KlinesRepository) GetAll(startTimestamp, endTimestamp int64) ([]Kline, int, error) {
	var klines []Kline

	query := r.db.Model(&Kline{})

	fmt.Println(startTimestamp)
	fmt.Println(endTimestamp)

	if startTimestamp > 0 {
		query = query.Where("open_time >= ?", startTimestamp)
	}

	if endTimestamp > 0 {
		query = query.Where("open_time <= ?", endTimestamp)
	}

	if err := query.Order("open_time DESC").Find(&klines).Error; err != nil {
		return nil, 0, err
	}
	return klines, len(klines), nil
}
