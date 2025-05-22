package data

import "time"

type Book struct {
	ID        int64
	CreatedAt time.Time
	Title     string
	Author    string
	Year      int32
	Size      int32
	Genres    []string
	Version   int32
}
