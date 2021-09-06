package models

import "time"

type Bookmarks struct {
	ID          string    `json:"id"`
	Bookmarks   string    `json:"bookmarks"`
	LastUpdated time.Time `json:"lastUpdated"`
	Version     string    `json:"version"`
}

func (b Bookmarks) Empty() bool {
	return b.ID == "" || b.LastUpdated.Equal(time.UnixMicro(0)) || b.Version == ""
}
