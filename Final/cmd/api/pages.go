package main

import (
	"Final/internal/data"
	"Final/internal/validator"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmpl *template.Template

type Profile struct {
	ID      int64            `json:"id"`
	Name    string           `json:"name"`
	Email   string           `json:"email"`
	Role    string           `json:"role"`
	Balance int64            `json:"balance"`
	Items   []data.CaseItems `json:"items"`
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	logout, _ := r.Cookie("token")
	logout.MaxAge = -1
	http.SetCookie(w, logout)
	http.Redirect(w, r, "/", 303)
}
func (app *application) InventoryHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Cookie("token")
	user, _ := app.models.Users.GetForToken(data.ScopeAuthentication, token.Value)
	inventory, _ := app.models.Inventory.GetByID(user.ID)

	tmpl = template.Must(template.ParseFiles("cmd/api/templates/inventory.html"))
	tmpl.ExecuteTemplate(w, "inventory.html", inventory)
}
func (app *application) storePageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Itemname string
		Type     string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Itemname = app.readString(qs, "name", "")
	input.Type = app.readString(qs, "type", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "type", "stars", "-id", "-itemname", "-type", "-stars"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Accept the metadata struct as a return value.
	items, _, err := app.models.Case.GetAllCase(input.Itemname, input.Type, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tmpl = template.Must(template.ParseFiles("cmd/api/templates/store.html"))
	tmpl.ExecuteTemplate(w, "store.html", items)
}
func (app *application) mainPageHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Cookie("token")
	if token != nil {
		user, _ := app.models.Users.GetForToken(data.ScopeAuthentication, token.Value)
		if user.Role == "admin" {
			tmpl = template.Must(template.ParseFiles("cmd/api/templates/adminindex.html"))
			tmpl.ExecuteTemplate(w, "adminindex.html", nil)
		}
		tmpl = template.Must(template.ParseFiles("cmd/api/templates/indexxLog.html"))
		tmpl.ExecuteTemplate(w, "indexxLog.html", nil)
	} else {
		tmpl = template.Must(template.ParseFiles("cmd/api/templates/indexx.html"))
		tmpl.ExecuteTemplate(w, "indexx.html", nil)
	}

}
func (app *application) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/login.html"))
	tmpl.ExecuteTemplate(w, "login.html", nil)
}
func (app *application) registerPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/register.html"))
	tmpl.ExecuteTemplate(w, "register.html", nil)
}
func (app *application) adminPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/adminpage.html"))
	tmpl.ExecuteTemplate(w, "adminpage.html", nil)
}
func (app *application) errorPageHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	error1 := app.readString(qs, "err", "")
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/errorMessage.html"))
	tmpl.ExecuteTemplate(w, "errorMessage.html", error1)
}
func (app *application) confirmPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/confirm.html", "cmd/api/templates/indexxLog.html"))
	tmpl.ExecuteTemplate(w, "confirm.html", nil)
}
func (app *application) permissionPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/permission.html"))
	tmpl.ExecuteTemplate(w, "permission.html", nil)
}
func (app *application) adminHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/admin.html"))
	tmpl.ExecuteTemplate(w, "admin.html", nil)
}
func (app *application) inputCaseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Itemname string
		Type     string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Itemname = app.readString(qs, "itemname", "")
	input.Type = app.readString(qs, "type", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "itemname", "type", "stars", "-id", "-itemname", "-type", "-stars"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	items, _, err := app.models.CaseItem.GetAllCaseItem(input.Itemname, input.Type, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/caseinput.html"))
	tmpl.ExecuteTemplate(w, "caseinput.html", items)
}
func (app *application) viewCaseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	x, _ := app.models.Case.GetCaseID(num)

	tmpl = template.Must(template.ParseFiles("cmd/api/templates/case.gohtml"))
	tmpl.ExecuteTemplate(w, "case.gohtml", x)
}
func (app *application) viewAllCaseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Itemname string
		Type     string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Itemname = app.readString(qs, "name", "")
	input.Type = app.readString(qs, "type", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "type", "stars", "-id", "-itemname", "-type", "-stars"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	items, _, err := app.models.Case.GetAllCase(input.Itemname, input.Type, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tmpl = template.Must(template.ParseFiles("cmd/api/templates/Allcases.html"))
	tmpl.ExecuteTemplate(w, "Allcases.html", items)
}
func (app *application) imageHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID of the image from the request URL
	vars := mux.Vars(r)
	id := vars["id"]
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var data1 []byte
	x, _ := app.models.CaseItem.GetCaseItem(num)
	data1 = x.Image
	w.Header().Set("Content-Type", "image/png") // Or "image/png" if the image is in PNG format

	// Write the image data to the HTTP response body
	_, err = w.Write(data1)
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) listCasesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Itemname string
		Type     string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Itemname = app.readString(qs, "itemname", "")
	input.Type = app.readString(qs, "type", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "itemname", "type", "stars", "-id", "-itemname", "-type", "-stars"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	items, _, err := app.models.CaseItem.GetAllCaseItem(input.Itemname, input.Type, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/listitems.html"))

	tmpl.ExecuteTemplate(w, "listitems.html", items)
}
func (app *application) profilePageHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Cookie("token")
	userFound, _ := app.models.Users.GetForToken(data.ScopeAuthentication, token.Value)
	inventory, _ := app.models.Inventory.GetByID(userFound.ID)
	profile := Profile{
		ID:      userFound.ID,
		Name:    userFound.Name,
		Email:   userFound.Email,
		Role:    userFound.Role,
		Balance: userFound.Balance,
		Items:   inventory.Items,
	}
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/profile.html"))
	tmpl.ExecuteTemplate(w, "profile.html", profile)
}
