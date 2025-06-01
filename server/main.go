package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"imageproxy/internal/cache"
	"imageproxy/internal/processor"
	Storage "imageproxy/internal/storage"
)

var (
	CacheCapacity int
	ImgStorage    Storage.Storage
)

func RunServer(cacheCapacity int) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	cache := cache.NewLRUCache(CacheCapacity, ImgStorage)
	processor := processor.NewImageProcessor(cache)

	// Хендлер для тестирования.
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/fill/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		width, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Invalid width", http.StatusBadRequest)
			return
		}

		height, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "Invalid height", http.StatusBadRequest)
			return
		}

		url := strings.Join(parts[4:], "/")
		if url == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		data, contentType, err := processor.ProcessImage(r.Context(), url, width, height)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			fmt.Printf("Failed to write response: %v\n", err)
		}
	})

	fmt.Printf("Server listening on :%s (cache capacity: %d)\n", port, cacheCapacity)
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,   // максимальное время чтения запроса
		WriteTimeout: 10 * time.Second,  // максимальное время записи ответа
		IdleTimeout:  120 * time.Second, // максимальное время ожидания следующего запроса
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

func cacheCapacity() int {
	cacheCapacity := 5
	if envCap := os.Getenv("CACHE_CAPACITY"); envCap != "" {
		if capacity, err := strconv.Atoi(envCap); err == nil {
			cacheCapacity = capacity
		}
	}
	return cacheCapacity
}

func main() {
	CacheCapacity = cacheCapacity()
	var err error
	if os.Getenv("STORAGE_TYPE") == "memory" {
		ImgStorage = Storage.NewMemoryStorage()
	} else {
		ImgStorage, err = Storage.NewFileStorage("./image_cache")
		if err != nil {
			fmt.Printf("Failed to initialize file ImgStoragetorage: %v\n", err)
			os.Exit(1)
		}
	}

	RunServer(cacheCapacity())
}
