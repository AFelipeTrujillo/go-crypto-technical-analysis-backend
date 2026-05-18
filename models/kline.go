package models

type Kline struct {
	OpenTime                 int64
	OpenPrice                float64
	HighPrice                float64
	LowPrice                 float64
	ClosePrice               float64
	Volume                   float64
	CloseTime                int64
	QuoteAssetVolume         float64
	NumberOfTrades           int64
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
	IgnoreField              float64
}

func (p *Kline) TableName() string {
	return "btc_candles_1h"
}
