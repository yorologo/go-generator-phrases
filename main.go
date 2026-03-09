package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/yorologo/go-generator-phrases/image"
)

const defaultBatchSize = 10

func imageHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func loadMoreImagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count, err := parseCount(r.URL.Query().Get("count"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imagePaths, err := image.GenerateImages(count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(imagePaths); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseCount(raw string) (int, error) {
	if raw == "" {
		return defaultBatchSize, nil
	}

	count, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errors.New("count must be a valid integer")
	}
	if count <= 0 || count > 50 {
		return 0, errors.New("count must be between 1 and 50")
	}

	return count, nil
}

func main() {
	if err := os.MkdirAll("img", 0o755); err != nil {
		log.Fatalf("create output directory: %v", err)
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	mux.HandleFunc("/", imageHandler(tmpl))
	mux.HandleFunc("/load-more-images", loadMoreImagesHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
