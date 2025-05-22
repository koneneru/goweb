package data

import "time"

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int32     `json:"year"`
	Size      int32     `json:"size"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
}
