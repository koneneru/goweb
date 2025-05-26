package data

import (
	"context"
	"database/sql"
	"errors"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&b.ID, &b.CreatedAt, &b.Version)
}

func (m BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, author, year, size, genres, version
		FROM books
		WHERE id=$1`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Author,
		&book.Year,
		&book.Size,
		pq.Array(&book.Genres),
		&book.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}

	return &book, nil
}

func (m BookModel) GetAll(title, author string, genres []string, filters Filters) ([]*Book, error) {
	query := `
		SELECT id,created_at,title,author,year,size,genres,version
		FROM books
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1='')
		AND (LOWER(author)=LOWER($2) OR $2='')
		AND (genres @> $3 OR $3='{}')
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, title, author, pq.Array(genres))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*Book{}

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&book.ID,
			&book.CreatedAt,
			&book.Title,
			&book.Author,
			&book.Year,
			&book.Size,
			pq.Array(&book.Genres),
			&book.Version,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (m BookModel) Update(b *Book) error {
	query := `
		UPDATE books
		SET title=$1, author=$2, year=$3, size=$4, genres=$5, version=version+1
		WHERE id=$6 AND version=$7
		RETURNING version`

	args := []any{
		b.Title,
		b.Author,
		b.Year,
		b.Size,
		pq.Array(b.Genres),
		b.ID,
		b.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&b.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict

		default:
			return err
		}
	}

	return nil
}

func (m BookModel) Delete(id int64) error {
	query := `
		DELETE FROM books
		WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
