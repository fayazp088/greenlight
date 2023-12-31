package data

import "time"

type Genre []string

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime"`
	Genres    Genre     `json:"genres"`
	Version   int32     `json:"version"`
}
