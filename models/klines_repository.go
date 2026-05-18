package models

import "gorm.io/gorm"

type KlinesRepositoryInterface interface {
	GetAll() ([]Kline, int, error)
}

type KlinesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) KlinesRepositoryInterface {
	return &KlinesRepository{
		db: db,
	}
}

func (r *KlinesRepository) GetAll() ([]Kline, int, error) {
	var klines []Kline
	if err := r.db.Order("open_time DESC").Find(&klines).Error; err != nil {
		return nil, 0, err
	}
	return klines, len(klines), nil
}
