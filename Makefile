tidy ::
	@go mod tidy && go mod vendor

seed ::
	@go run cmd/seed/main.go

binance ::
	@go run cmd/binance/main.go

run ::
	@go run cmd/server/main.go