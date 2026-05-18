CREATE TABLE IF NOT EXISTS btc_candles_1h (
    open_time INTEGER PRIMARY KEY, 
    open_price REAL NOT NULL,
    high_price REAL NOT NULL,
    low_price REAL NOT NULL,
    close_price REAL NOT NULL,
    volume REAL NOT NULL,
    close_time INTEGER NOT NULL,
    quote_asset_volume REAL NOT NULL,
    number_of_trades INTEGER NOT NULL,
    taker_buy_base_asset_volume REAL NOT NULL,
    taker_buy_quote_asset_volume REAL NOT NULL,
    ignore_field REAL
);

-- Crear el índice en SQLite si no existe
CREATE INDEX IF NOT EXISTS idx_btc_candles_time ON btc_candles_1h (open_time DESC);