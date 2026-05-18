package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AFelipeTrujillo/go-crypto-technical-analysis-backend/app/database"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db, close := database.New()
	defer close()
	fmt.Println("Conected to DB...")

	dir := os.Getenv("SQLITE_DIR")
	fmt.Println("Reading " + dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("reading directory failed: %v", err)
	}

	var sqlFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
	}

	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	for _, file := range sqlFiles {
		filePath := filepath.Join(dir, file.Name())
		fmt.Printf("→ Execute: %s... \n", file.Name())

		query, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("\nError al leer el archivo %s: %v", file.Name(), err)
		}

		if err := db.Exec(string(query)).Error; err != nil {
			log.Printf("executing %s failed: %v", file.Name(), err)
			return
		}

		fmt.Println("→ OK")
	}
}
