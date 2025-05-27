package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	mux.HandleFunc("GET /v1/books", app.listBooksHandler)
	mux.HandleFunc("POST /v1/books", app.createBookHandler)
	mux.HandleFunc("GET /v1/books/{id}", app.showBookHandler)
	mux.HandleFunc("PATCH /v1/books/{id}", app.updateBookHandler)
	mux.HandleFunc("DELETE /v1/books/{id}", app.deleteBookHandler)

	mux.HandleFunc("POST /v1/users", app.registerUserHandler)

	return app.recoverPanic(app.rateLimit(mux))
}
