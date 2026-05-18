package klines

import (
	"net/http"

	"github.com/AFelipeTrujillo/go-crypto-technical-analysis-backend/app/api"
	"github.com/AFelipeTrujillo/go-crypto-technical-analysis-backend/models"
)

type KlineResponse struct {
	Klines []Kline `json:"klines"`
	Total  int     `json:"total"`
}

type Kline struct {
	OpenTime   int64   `json:"open_time"`
	OpenPrice  float64 `json:"open_price"`
	HighPrice  float64 `json:"high_price"`
	LowPrice   float64 `json:"low_price"`
	ClosePrice float64 `json:"close_price"`
	Volume     float64 `json:"volume"`
	CloseTime  int64   `json:"close_time"`
}

type KlineHandler struct {
	repo models.KlinesRepositoryInterface
}

func NewKlinesHandler(r models.KlinesRepositoryInterface) *KlineHandler {
	return &KlineHandler{
		repo: r,
	}
}

func (h *KlineHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {

	res, total, err := h.repo.GetAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	klines := make([]Kline, len(res))
	for i, k := range res {
		klines[i] = Kline{
			OpenTime:   k.OpenTime,
			OpenPrice:  k.OpenPrice,
			HighPrice:  k.HighPrice,
			LowPrice:   k.LowPrice,
			ClosePrice: k.ClosePrice,
			Volume:     k.Volume,
			CloseTime:  k.CloseTime,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := KlineResponse{
		Klines: klines,
		Total:  total,
	}

	api.OKResponse(w, response)
}
