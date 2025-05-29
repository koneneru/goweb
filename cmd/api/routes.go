package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	mux.HandleFunc("GET /v1/books", app.requirePermission("books:read", app.listBooksHandler))
	mux.HandleFunc("POST /v1/books", app.requirePermission("books:write", app.createBookHandler))
	mux.HandleFunc("GET /v1/books/{id}", app.requirePermission("books:read", app.showBookHandler))
	mux.HandleFunc("PATCH /v1/books/{id}", app.requirePermission("books:write", app.updateBookHandler))
	mux.HandleFunc("DELETE /v1/books/{id}", app.requirePermission("books:write", app.deleteBookHandler))

	mux.HandleFunc("POST /v1/users", app.registerUserHandler)
	mux.HandleFunc("PUT /v1/users/activate", app.activateUserHandler)

	mux.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(mux))))
}
