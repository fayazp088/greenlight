package models

import (
	"time"

	"github.com/fayazp088/greenlight/internal/data"
)

type Movie struct {
	ID        int64        `json:"id"`
	CreatedAt time.Time    `json:"-"`
	Title     string       `json:"title"`
	Year      int32        `json:"year,omitempty"`
	Runtime   data.Runtime `json:"runtime,omitempty"`
	Genres    []string     `json:"genres,omitempty"`
	Version   int32        `json:"version"`
}
