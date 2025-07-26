package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func tableCheck() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_CONTAINER_PORT"), os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_SSLMODE"),
	)

	fmt.Printf("DB URL: %s", dsn)

	db, err := sql.Open(os.Getenv("POSTGRES_DRIVER"), dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		);
	`
	var tableExists bool
	err = db.QueryRowContext(context.Background(), query, "stocks").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("table does not exists: %w", err)
	}

	return nil
}

func dbHealthCheck(w http.ResponseWriter, r *http.Request) {
	err := tableCheck()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "db check error: %v", err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "db access and migration check completed")
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "http server works correctly")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthCheck)
	mux.HandleFunc("GET /db", dbHealthCheck)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), mux))
}
