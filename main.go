package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HTTP server works correctly")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthCheck)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), mux))
}
