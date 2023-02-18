package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("./cmd"))
	router.PathPrefix("/cmd/").Handler(http.StripPrefix("/cmd/", fs))
	router.HandleFunc("/", app.mainPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/reg", app.RegPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/case", app.casePageHandler).Methods("GET", "OPTIONS")
	return app.recoverPanic(app.rateLimit(router))

}
