package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/tidwall/buntdb"

	"github.com/thinkofher/lalyta/pkg/api"
	"github.com/thinkofher/lalyta/pkg/storage"
)

func run() error {
	bunt, err := buntdb.Open("lalyta.db")
	if err != nil {
		return fmt.Errorf("buntdb.Open: %w", err)
	}
	defer bunt.Close()

	buntStorage := storage.New(bunt)

	mux := pat.New()
	mux.Get("/info", api.Info("PL", "Hello World!", "1.1.13"))
	mux.Post("/bookmarks", api.CreateBookmarks(buntStorage))
	mux.Get("/bookmarks/:id", api.Bookmarks(buntStorage))
	mux.Put("/bookmarks/:id", api.UpdateBookmarks(buntStorage))
	mux.Get("/bookmarks/:id/lastUpdated", api.LastUpdated(buntStorage))
	mux.Get("/bookmarks/:id/version", api.Version(buntStorage))

	log.Println("Starting server at 0.0.0.0:8080")
	return http.ListenAndServe("0.0.0.0:8080", mux)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
