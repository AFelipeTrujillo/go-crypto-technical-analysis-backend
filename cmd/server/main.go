package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AFelipeTrujillo/go-crypto-technical-analysis-backend/app/api/klines"
	"github.com/AFelipeTrujillo/go-crypto-technical-analysis-backend/app/database"
	"github.com/AFelipeTrujillo/go-crypto-technical-analysis-backend/models"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, close := database.New()
	defer close()
	fmt.Println("Conected to DB...")

	klinesRepo := models.NewCategoriesRepository(db)
	klinesHandler := klines.NewKlinesHandler(klinesRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/klines", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			klinesHandler.HandleGetAll(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}

		log.Println("Server stopped gracefully")
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	srv.Shutdown(ctx)
	stop()

}
