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
	//auth
	router.HandleFunc("/login", app.loginPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", app.loginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/register", app.registerUserHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/register", app.registerPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/logout", app.logoutHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/delete/{id}", app.requireAuthenticatedUser(app.deleteUserHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/profile", app.requireAuthenticatedUser(app.updateUserHandler)).Methods("POST", "OPTIONS")
	//confirm
	router.HandleFunc("/v1/users/activated", app.activateUserHandler).Methods("GET", "OPTIONS")
	//messages
	router.HandleFunc("/confirm", app.confirmPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/permission", app.permissionPageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/errors", app.errorPageHandler).Methods("GET", "OPTIONS")

	//crud
	router.HandleFunc("/admin", app.requireAdminUser(app.adminPageHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/adminpage", app.requireAdminUser(app.adminHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/profile", app.requireAuthenticatedUser(app.profilePageHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/list", app.requireAdminUser(app.listCasesHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/admin", app.requireAdminUser(app.createCaseItemsHandler)).Methods("POST", "OPTIONS")

	router.HandleFunc("/caseadd", app.requireAdminUser(app.inputCaseHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/caseadd", app.requireAdminUser(app.createCaseHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/case/{id}", app.viewCaseHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/case", app.requireAdminUser(app.viewAllCaseHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/deletecase/{id}", app.requireAdminUser(app.deleteCaseHandler)).Methods("GET", "OPTIONS")

	router.HandleFunc("/getitem/{id}", app.requireAuthenticatedUser(app.toInventoryHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/inventory", app.InventoryHandler).Methods("GET", "OPTIONS")

	router.HandleFunc("/store", app.storePageHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/buy/{id}", app.balanceMinusHandler).Methods("GET", "OPTIONS")

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))

}
