package models

import "time"

type Movie struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int       `json:"year"`
	Genres    []string  `json:"genres"`
	Version   int       `json:"version"`
}
