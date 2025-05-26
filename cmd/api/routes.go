package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	mux.HandleFunc("POST /v1/books", app.createBookHandler)
	mux.HandleFunc("GET /v1/books/{id}", app.showBookHandler)
	mux.HandleFunc("PUT /v1/books/{id}", app.updateBookHandler)
	mux.HandleFunc("DELETE /v1/books/{id}", app.deleteBookHandler)

	return mux
}
