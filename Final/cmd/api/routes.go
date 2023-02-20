package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/images/{id}", app.imageHandler).Methods("GET", "OPTIONS")
	fs := http.FileServer(http.Dir("./cmd"))
	router.PathPrefix("/cmd/").Handler(http.StripPrefix("/cmd/", fs))
	router.HandleFunc("/", app.mainPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", app.loginPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", app.loginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/register", app.registerUserHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/register", app.registerPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/logout", app.logoutHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/v1/users/activated", app.activateUserHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/confirm", app.confirmPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/permission", app.permissionPageHandler).Methods("GET", "OPTIONS")

	//crud
	router.HandleFunc("/adminpage", app.requireAdminUser(app.adminPageHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/list", app.requireAdminUser(app.listCasesHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/adminpage", app.requireAdminUser(app.createCaseItemsHandler)).Methods("POST", "OPTIONS")

	router.HandleFunc("/caseadd", app.inputCaseHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/caseadd", app.createCaseHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/case/{id}", app.viewCaseHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/case", app.viewAllCaseHandler).Methods("GET", "OPTIONS")

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))

}
