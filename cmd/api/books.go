package main

import (
	"fmt"
	"goweb/internal/data"
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
