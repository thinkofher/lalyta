package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/thinkofher/lalyta/pkg/api"
	"github.com/thinkofher/lalyta/pkg/models"
	"github.com/tidwall/buntdb"
)

type DB struct {
	bunt *buntdb.DB
}

func New(b *buntdb.DB) *DB {
	b.CreateIndex("bookmarks", "bookmarks:*", buntdb.IndexJSON("id"))
	return &DB{
		bunt: b,
	}
}

type bookmarksEntry struct {
	ID          string    `json:"id"`
	Bookmarks   string    `json:"bookmarks"`
	LastUpdated time.Time `json:"lastUpdated"`
	Version     string    `json:"version"`
}

func bookmarksKey(id string) string {
	return fmt.Sprintf("bookmarks:%s", id)
}

func makeBookmarksEntry(b bookmarksEntry) (string, string) {
	val, err := json.Marshal(b)
	if err != nil {
		return "", ""
	}
	return bookmarksKey(b.ID), string(val)
}

func (db *DB) SetBookmarks(ctx context.Context, b models.Bookmarks) error {
	return db.bunt.Update(func(tx *buntdb.Tx) error {
		key, value := makeBookmarksEntry(bookmarksEntry{
			ID:          b.ID,
			Bookmarks:   b.Bookmarks,
			LastUpdated: b.LastUpdated,
			Version:     b.Version,
		})

		_, _, err := tx.Set(key, value, nil)
		if err != nil {
			return fmt.Errorf("tx.Set: %w", err)
		}
		return nil
	})
}

func (db *DB) GetBookmarks(ctx context.Context, id string) (*models.Bookmarks, error) {
	res := new(models.Bookmarks)

	err := db.bunt.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(bookmarksKey(id))
		if err != nil {
			return fmt.Errorf("tx.Get: %w", err)
		}

		b := new(bookmarksEntry)

		err = json.Unmarshal([]byte(val), b)
		if err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}

		res = &models.Bookmarks{
			ID:          b.ID,
			Bookmarks:   b.Bookmarks,
			LastUpdated: b.LastUpdated,
			Version:     b.Version,
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("db.bunt.View: %w", err)
	}
	if res.Empty() {
		return nil, api.ErrBookmarksNotFound
	}

	return res, nil
}
