/*
Package api implements http REST API compatible with xBrowserSync service.

Below handlers contains information about what method should they be used with
and what path should they be mounted to. It is dictated by xBrowserSync original
API implementation.

Handlers implemented by api package are completly indepedent of used http router.
It is up to implementing person to choose proper http router.

See xBrowserSync documentation for further information about what it is and how it
works.
*/
package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alioygur/gores"

	"github.com/thinkofher/lalyta/pkg/models"
	"github.com/thinkofher/lalyta/pkg/service/gen"
)

// Info retrieves information describing the xBrowserSync service.
//
//  GET /info
//
// Response example:
//
//  {
//    "maxSyncSize": 204800,
//    "message": "",
//    "status": 1,
//    "version": "1.1.13"
//  }
//
// * Status ("status", int): current service status code. 1 = Online; 2 = Offline;
// 3 = Not accepting new syncs.
//
// * Message ("message", string): service information message.
//
// * Version ("version", string): API version service is using.
//
// * Maximum sync size ("maxSyncSize", int): maximum sync size (in bytes)
// allowed by the service.
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

// BookmarksStorage describes methods required for storing and retrieving
// bookmarks from data source. It has to be thread safe.
type BookmarksStorage interface {
	// SetBookmarks method stores given encrypted Bookmarks
	// in database.
	SetBookmarks(ctx context.Context, b models.Bookmarks) error

	// GetBookmarks retrieves encrypted Bookmarks with given id
	// from database.
	GetBookmarks(ctx context.Context, id string) (*models.Bookmarks, error)
}

// QueryParameters help to determine specific bookmarks.
type QueryParameters interface {
	// ID returns 32 character alphanumeric sync ID used
	// by client.
	ID(*http.Request) string
}

// ErrBookmarksNotFound is returned by BookmarksStorage when
// there is no bookmarks with given ID in storage.
var ErrBookmarksNotFound = errors.New("bookmarks with given id has been not found")

// CreateBookmarks creates a new (empty) bookmark sync and returns
// the corresponding ID.
//
//  POST /bookmarks
//
// Post body example:
//
//  {
//    "version": "1.0.0"
//  }
//
// Response example:
//
//  {
//    "id": "52758cb942814faa9ab255208025ae59",
//    "lastUpdated": "2016-07-06T12:43:16.866Z",
//    "version": "1.0.0"
//  }
//
// * ID ("id", string): 32 character alphanumeric sync ID.
//
// * Last updated ("lastUpdated", timestamp as string): last updated timestamp
// for created bookmarks.
//
// * Version ("version", version): version number of the xBrowserSync client
// used to create the sync.
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

// Bookmarks retrieves the bookmark sync corresponding to the
// provided sync ID.
//
//  GET /bookmarks/{id}
//
// Query params:
//
// * id: 32 character alphanumeric sync ID.
//
// Response example:
//
//  {
//    "bookmarks": "DWCx6wR9ggPqPRrhU4O4oLN5P09oULX4Xt+ckxswtFNds...",
//    "lastUpdated": "2016-07-06T12:43:16.866Z",
//    "version": "1.0.0"
//  }
//
// * Bookmarks ("bookmarks", string): encrypted bookmark data salted using
// secret value.
//
// * Last updated: ("lastUpdated", timestamp as string): last updated timestamp
// for retrieved bookmarks.
//
// * Version ("version", string): version number of the xBrowserSync client used
// to create the sync.
func Bookmarks(storage BookmarksStorage, params QueryParameters) http.HandlerFunc {
	type response struct {
		Bookmarks   string    `json:"bookmarks"`
		LastUpdated time.Time `json:"lastUpdated"`
		Version     string    `json:"version"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := params.ID(r)
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

// UpdateBookmarks updates the bookmark sync data corresponding to the
// provided sync ID with the provided encrypted bookmarks data.
//
//  PUT /bookmarks/{id}
//
// Query params:
//
// * id: 32 character alphanumeric sync ID.
//
// Post body example:
//
//  {
//    "bookmarks": "DWCx6wR9ggPqPRrhU4O4oLN5P09oULX4Xt+ckxswtFNds...",
//    "lastUpdated": "2016-07-06T12:43:16.866Z",
//  }
//
// * Bookmarks ("bookmarks", string): encrypted bookmark data salted using
// secret value.
//
// * Last updated ("lastUpdated", timestamp as string): last updated timestamp
// to check against existing bookmarks.
//
// Response example:
//
//  {
//    "lastUpdated": "2016-07-06T12:43:16.866Z"
//  }
//
// Last updated ("lastUpdated", timestamp as string): last updated timestamp
// for updated bookmarks.
func UpdateBookmarks(storage BookmarksStorage, params QueryParameters) http.HandlerFunc {
	type payload struct {
		Bookmarks   string    `json:"bookmarks"`
		LastUpdated time.Time `json:"lastUpdated"`
	}
	type response struct {
		LastUpdated time.Time `json:"lastUpdated"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := params.ID(r)
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

// LastUpdated retrieves the bookmark sync last updated timestamp
// corresponding to the provided sync ID.
//
//   GET /bookmarks/{id}/lastUpdated
//
// Query params:
//
// * id: 32 character alphanumeric sync ID.
//
// Response example:
//
//  {
//    "lastUpdated":"2016-07-06T12:43:16.866Z"
//  }
//
// * Last updated ("lastUpdated", timestamp as string): last updated
// timestamp for corresponding bookmarks.
func LastUpdated(storage BookmarksStorage, params QueryParameters) http.HandlerFunc {
	type response struct {
		LastUpdated time.Time `json:"lastUpdated"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := params.ID(r)
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

// Version retrieves the bookmark sync version number of the xBrowserSync client
// used to create the bookmarks sync corresponding to the provided sync ID.
//
//  GET /bookmarks/{id}/version
//
// Query params:
//
// * id: 32 character alphanumeric sync ID.
//
// Response example:
//
//  {
//    "version":"1.0.0"
//  }
//
// Version ("version", string): version number of the xBrowserSync client
// used to create the sync.
func Version(storage BookmarksStorage, params QueryParameters) http.HandlerFunc {
	type response struct {
		Version string `json:"version"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := params.ID(r)
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
