-- 1. Create symbols table
CREATE TABLE IF NOT EXISTS symbols (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Insert initial symbol (BTCUSDT)
INSERT OR IGNORE INTO symbols (name) VALUES ('BTCUSDT');

-- 3. Since SQLite doesn't support ADD COLUMN with PK or FK constraints easily,
-- we need to recreate klines_1h with the new column and constraints.

-- 3.1 Create the new table with the desired schema
CREATE TABLE IF NOT EXISTS klines_1h_new (
    symbol_id INTEGER NOT NULL,
    open_time INTEGER NOT NULL,
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
    ignore_field REAL,
    PRIMARY KEY (symbol_id, open_time),
    FOREIGN KEY (symbol_id) REFERENCES symbols(id)
);

-- 3.2 Copy data from old table assigning symbol_id = 1 (BTCUSDT)
INSERT INTO klines_1h_new (
    symbol_id, open_time, open_price, high_price, low_price, close_price, volume,
    close_time, quote_asset_volume, number_of_trades,
    taker_buy_base_asset_volume, taker_buy_quote_asset_volume, ignore_field
)
SELECT
    1, open_time, open_price, high_price, low_price, close_price, volume,
    close_time, quote_asset_volume, number_of_trades,
    taker_buy_base_asset_volume, taker_buy_quote_asset_volume, ignore_field
FROM klines_1h;

-- 3.3 Drop old table and rename new one
DROP TABLE klines_1h;
ALTER TABLE klines_1h_new RENAME TO klines_1h;

-- 4. Create indexes
CREATE INDEX IF NOT EXISTS idx_klines_1h_symbol_time ON klines_1h (symbol_id, open_time DESC);
CREATE INDEX IF NOT EXISTS idx_klines_1h_open_time ON klines_1h (open_time DESC);
