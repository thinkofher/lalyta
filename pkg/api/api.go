package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alioygur/gores"
	"github.com/go-chi/chi/v5"
	"github.com/thinkofher/lalyta/pkg/models"
	"github.com/thinkofher/lalyta/pkg/service/gen"
)

func Info(location, msg, version string) http.HandlerFunc {
	type response struct {
		MaxSyncSize int64  `json:"maxSyncSize"`
		Message     string `json:"message"`
		Status      int    `json:"status"`
		Version     string `json:"version"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gores.JSON(w, http.StatusOK, &response{
			MaxSyncSize: 204800,
			Message:     msg,
			Status:      1,
			Version:     version,
		})
	})
}

const idLength = 32

type BookmarksStorage interface {
	SetBookmarks(ctx context.Context, b models.Bookmarks) error
	GetBookmarks(ctx context.Context, id string) (*models.Bookmarks, error)
}

var ErrBookmarksNotFound = errors.New("bookmarks with given id has been not found")

func CreateBookmarks(storage BookmarksStorage) http.HandlerFunc {
	type payload struct {
		Version string `json:"version"`
	}
	type response struct {
		ID          string    `json:"id"`
		LastUpdated time.Time `json:"lastUpdated"`
		Version     string    `json:"version"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		p := new(payload)

		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		id, err := gen.String(idLength)
		if err != nil {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bookmarks := models.Bookmarks{
			ID:          id,
			Bookmarks:   "",
			LastUpdated: time.Now().UTC(),
			Version:     p.Version,
		}
		if err := storage.SetBookmarks(ctx, bookmarks); err != nil {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		gores.JSON(w, http.StatusOK, &response{
			ID:          bookmarks.ID,
			LastUpdated: bookmarks.LastUpdated,
			Version:     bookmarks.Version,
		})
	})
}

func Bookmarks(storage BookmarksStorage) http.HandlerFunc {
	type response struct {
		Bookmarks   string    `json:"bookmarks"`
		LastUpdated time.Time `json:"lastUpdated"`
		Version     string    `json:"version"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		if id == "" {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := storage.GetBookmarks(ctx, id)
		if err != nil {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusNotFound)
			return
		}

		gores.JSON(w, http.StatusOK, &response{
			Bookmarks:   b.Bookmarks,
			LastUpdated: b.LastUpdated,
			Version:     b.Version,
		})
	})
}

func UpdateBookmarks(storage BookmarksStorage) http.HandlerFunc {
	type payload struct {
		Bookmarks   string    `json:"bookmarks"`
		LastUpdated time.Time `json:"lastUpdated"`
	}
	type response struct {
		LastUpdated time.Time `json:"lastUpdated"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		if id == "" {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p := new(payload)
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := storage.GetBookmarks(ctx, id)
		if err != nil {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if p.LastUpdated != b.LastUpdated {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		now := time.Now().UTC()
		err = storage.SetBookmarks(ctx, models.Bookmarks{
			ID:          id,
			Bookmarks:   p.Bookmarks,
			LastUpdated: now,
			Version:     b.Version,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		gores.JSON(w, http.StatusOK, &response{
			LastUpdated: now,
		})
	})
}

func LastUpdated(storage BookmarksStorage) http.HandlerFunc {
	type response struct {
		LastUpdated time.Time `json:"lastUpdated"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		if id == "" {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := storage.GetBookmarks(ctx, id)
		if err != nil {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusNotFound)
			return
		}

		gores.JSON(w, http.StatusOK, &response{
			LastUpdated: b.LastUpdated,
		})
	})
}

func Version(storage BookmarksStorage) http.HandlerFunc {
	type response struct {
		Version string `json:"version"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		if id == "" {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := storage.GetBookmarks(ctx, id)
		if err != nil {
			// TODO(thinkofher) output json error message
			w.WriteHeader(http.StatusNotFound)
			return
		}

		gores.JSON(w, http.StatusOK, &response{
			Version: b.Version,
		})
	})
}
