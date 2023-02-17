package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/api/items", app.listCasesHandler)
	router.HandlerFunc(http.MethodPost, "/api/items", app.createCaseItemsHandler)
	router.HandlerFunc(http.MethodGet, "/api/items/:id", app.showItemsHandler)
	router.HandlerFunc(http.MethodDelete, "/api/items/:id", app.deleteItemHandler)
	router.HandlerFunc(http.MethodPut, "/api/items/:id", app.updateItemCaseHandler)
	return app.recoverPanic(app.rateLimit(router))

}
