package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new book")
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("id")

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of book %d\n", id)
}
