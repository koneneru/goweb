package main

import (
	"fmt"
	"goweb/internal/data"
	"net/http"
	"time"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new book")
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
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
		app.logger.Println(err)
		http.Error(w, "The server encounterred a problem and could not process your request", http.StatusInternalServerError)
	}
}
