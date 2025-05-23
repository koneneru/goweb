package data

import (
	"database/sql"
	"goweb/internal/validator"
	"time"

	"github.com/lib/pq"
)

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int32     `json:"year,omitempty"`
	Size      Booksize  `json:"size,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 500, "title", "must not be greater than 500 bytes long")

	v.Check(book.Author != "", "author", "must be provided")
	v.Check(len(book.Author) <= 200, "author", "must not be greater than 100 bytes long")

	v.Check(book.Year != 0, "year", "must be provided")
	v.Check(book.Year >= 1888, "year", "must bo greater than 1888")
	v.Check(book.Year <= int32(time.Now().Year()), "year", "must not be in future")

	v.Check(book.Size != 0, "size", "must be provided")
	v.Check(book.Size > 0, "size", "must be positive integer")

	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(book.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")
}

type BookModel struct {
	DB *sql.DB
}

func (m BookModel) Insert(b *Book) error {
	query := `
		INSERT INTO books (title, author, year, size, genres)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`

	args := []any{b.Title, b.Author, b.Year, b.Size, pq.Array(b.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&b.ID, &b.CreatedAt, &b.Version)
}

func (m BookModel) Get(id int64) (*Book, error) {
	return nil, nil
}

func (m BookModel) Update(b *Book) error {
	return nil
}

func (m BookModel) Delete(id int64) error {
	return nil
}
