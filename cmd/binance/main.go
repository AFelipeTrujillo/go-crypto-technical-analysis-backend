package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/AFelipeTrujillo/go-crypto-technical-analysis-backend/app/database"
	"github.com/AFelipeTrujillo/go-crypto-technical-analysis-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// parseKline converts a raw []interface{} from Binance JSON into a typed Kline struct.
// Indexes based on Binance API documentation for klines.
func parseKline(raw []interface{}) (models.Kline, error) {
	var k models.Kline

	if v, ok := raw[0].(float64); ok {
		k.OpenTime = int64(v)
	} else {
		log.Printf("Warning: open_time unexpected type: %T", raw[0])
	}

	if v, err := strconv.ParseFloat(raw[1].(string), 64); err == nil {
		k.OpenPrice = v
	} else {
		return k, err
	}
	if v, err := strconv.ParseFloat(raw[2].(string), 64); err == nil {
		k.HighPrice = v
	} else {
		return k, err
	}
	if v, err := strconv.ParseFloat(raw[3].(string), 64); err == nil {
		k.LowPrice = v
	} else {
		return k, err
	}
	if v, err := strconv.ParseFloat(raw[4].(string), 64); err == nil {
		k.ClosePrice = v
	} else {
		return k, err
	}
	if v, err := strconv.ParseFloat(raw[5].(string), 64); err == nil {
		k.Volume = v
	} else {
		return k, err
	}

	if v, ok := raw[6].(float64); ok {
		k.CloseTime = int64(v)
	} else {
		log.Printf("Warning: close_time unexpected type: %T", raw[6])
	}

	if v, err := strconv.ParseFloat(raw[7].(string), 64); err == nil {
		k.QuoteAssetVolume = v
	} else {
		return k, err
	}

	if v, ok := raw[8].(float64); ok {
		k.NumberOfTrades = int64(v)
	} else {
		log.Printf("Warning: number_of_trades unexpected type: %T", raw[8])
	}

	if v, err := strconv.ParseFloat(raw[9].(string), 64); err == nil {
		k.TakerBuyBaseAssetVolume = v
	} else {
		return k, err
	}
	if v, err := strconv.ParseFloat(raw[10].(string), 64); err == nil {
		k.TakerBuyQuoteAssetVolume = v
	} else {
		return k, err
	}
	if v, err := strconv.ParseFloat(raw[11].(string), 64); err == nil {
		k.IgnoreField = v
	} else {
		return k, err
	}

	return k, nil
}

// getOrCreateSymbol finds a symbol by name or creates it if it doesn't exist.
func getOrCreateSymbol(db *gorm.DB, name string) (int64, error) {
	var symbol models.Symbol
	result := db.Where("name = ?", name).First(&symbol)
	if result.Error == nil {
		return symbol.ID, nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, result.Error
	}

	symbol = models.Symbol{Name: name}
	if err := db.Create(&symbol).Error; err != nil {
		return 0, err
	}
	return symbol.ID, nil
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db, close := database.New()
	defer close()

	// Get symbol from env, default to BTCUSDT
	symbolName := os.Getenv("SYMBOL")
	if symbolName == "" {
		symbolName = "BTCUSDT"
	}

	// Get or create the symbol in the database
	symbolID, err := getOrCreateSymbol(db, symbolName)
	if err != nil {
		log.Fatalf("Error getting/creating symbol: %v", err)
	}

	binanceAPIURL := os.Getenv("BINANCE_API_BASE_URL")
	klinesAPIURL := binanceAPIURL + "/api/v3/klines?symbol=" + symbolName + "&interval=1h&limit=1000"

	resp, err := http.Get(klinesAPIURL)
	if err != nil {
		log.Fatalf("Error making Binance request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Binance returned error status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var rawKlines [][]interface{}
	if err := json.Unmarshal(body, &rawKlines); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Use GORM transaction
	tx := db.Begin()

	for _, raw := range rawKlines {
		if len(raw) < 12 {
			continue
		}

		kline, err := parseKline(raw)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Error parsing kline data: %v", err)
		}

		kline.SymbolID = symbolID

		result := tx.Exec(`
			INSERT OR IGNORE INTO klines_1h (
				symbol_id, open_time, open_price, high_price, low_price, close_price, volume,
				close_time, quote_asset_volume, number_of_trades,
				taker_buy_base_asset_volume, taker_buy_quote_asset_volume, ignore_field
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
		`,
			kline.SymbolID,
			kline.OpenTime,
			kline.OpenPrice,
			kline.HighPrice,
			kline.LowPrice,
			kline.ClosePrice,
			kline.Volume,
			kline.CloseTime,
			kline.QuoteAssetVolume,
			kline.NumberOfTrades,
			kline.TakerBuyBaseAssetVolume,
			kline.TakerBuyQuoteAssetVolume,
			kline.IgnoreField,
		)

		if result.Error != nil {
			tx.Rollback()
			log.Fatalf("Error inserting data: %v", result.Error)
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	log.Printf("Klines for %s inserted successfully", symbolName)
}
