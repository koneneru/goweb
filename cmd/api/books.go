package main

import (
	"fmt"
	"goweb/internal/data"
	"goweb/internal/validator"
	"net/http"
	"time"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string        `json:"title"`
		Author string        `json:"author"`
		Year   int32         `json:"year"`
		Size   data.Booksize `json:"size"`
		Genres []string      `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badReqestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be greater than 500 bytes long")

	v.Check(input.Author != "", "author", "must be provided")
	v.Check(len(input.Author) <= 200, "author", "must not be greater than 100 bytes long")

	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must bo greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in future")

	v.Check(input.Size != 0, "size", "must be provided")
	v.Check(input.Size > 0, "size", "must be positive integer")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book := data.Book{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Spice and Wolf vol.1",
		Author:    "Isuna Hasekura",
		Size:      289,
		Genres:    []string{"adventure", "drama", "fantasy", "romance", "supernatural"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
