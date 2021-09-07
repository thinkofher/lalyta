package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tidwall/buntdb"

	"github.com/thinkofher/lalyta/pkg/api"
	"github.com/thinkofher/lalyta/pkg/service/params"
	"github.com/thinkofher/lalyta/pkg/storage"
)

func run() error {
	bunt, err := buntdb.Open("lalyta.db")
	if err != nil {
		return fmt.Errorf("buntdb.Open: %w", err)
	}
	defer bunt.Close()

	buntStorage := storage.New(bunt)

	chiParams := new(params.Chi)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/info", api.Info("PL", "Hello World!", "1.1.13"))
	r.Post("/bookmarks", api.CreateBookmarks(buntStorage))
	r.Get("/bookmarks/{id}", api.Bookmarks(buntStorage, chiParams))
	r.Put("/bookmarks/{id}", api.UpdateBookmarks(buntStorage, chiParams))
	r.Get("/bookmarks/{id}/lastUpdated", api.LastUpdated(buntStorage, chiParams))
	r.Get("/bookmarks/{id}/version", api.Version(buntStorage, chiParams))

	log.Println("Starting server at 0.0.0.0:8080")
	return http.ListenAndServe("0.0.0.0:8080", r)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
