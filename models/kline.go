package models

type Symbol struct {
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"uniqueIndex;not null"`
}

type Kline struct {
	SymbolID                 int64   `gorm:"primaryKey;autoIncrement:false"`
	OpenTime                 int64   `gorm:"primaryKey;autoIncrement:false"`
	OpenPrice                float64 `gorm:"not null"`
	HighPrice                float64 `gorm:"not null"`
	LowPrice                 float64 `gorm:"not null"`
	ClosePrice               float64 `gorm:"not null"`
	Volume                   float64 `gorm:"not null"`
	CloseTime                int64   `gorm:"not null"`
	QuoteAssetVolume         float64 `gorm:"not null"`
	NumberOfTrades           int64   `gorm:"not null"`
	TakerBuyBaseAssetVolume  float64 `gorm:"not null"`
	TakerBuyQuoteAssetVolume float64 `gorm:"not null"`
	IgnoreField              float64

	Symbol Symbol `gorm:"foreignKey:SymbolID;references:ID"`
}

func (Kline) TableName() string {
	return "klines_1h"
}

func (Symbol) TableName() string {
	return "symbols"
}
