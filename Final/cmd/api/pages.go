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

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	logout, _ := r.Cookie("token")
	logout.MaxAge = -1
	http.SetCookie(w, logout)
}

func (app *application) mainPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/indexx.html"))
	tmpl.ExecuteTemplate(w, "indexx.html", nil)
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

func (app *application) confirmPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/confirm.html"))
	tmpl.ExecuteTemplate(w, "confirm.html", nil)
}
func (app *application) permissionPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/permission.html"))
	tmpl.ExecuteTemplate(w, "permission.html", nil)
}
func (app *application) inputCaseHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/caseinput.html"))
	tmpl.ExecuteTemplate(w, "caseinput.html", nil)
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
	// Accept the metadata struct as a return value.
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
	fmt.Println(id)
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var data1 []byte
	x, _ := app.models.CaseItem.GetCaseItem(num)
	fmt.Println(err)
	data1 = x.Image
	w.Header().Set("Content-Type", "image/png") // Or "image/png" if the image is in PNG format

	// Write the image data to the HTTP response body
	_, err = w.Write(data1)
	if err != nil {
		log.Fatal(err)
	}
}
func (app *application) casePageHandler(w http.ResponseWriter, r *http.Request) {
	cases := make([]data.CaseItems, 0)
	caseItem1 := data.CaseItems{
		ID:              0,
		ItemName:        "Diluc",
		ItemDescription: "Very sexy",
		Type:            "character",
		Stars:           5,
	}
	caseItem2 := data.CaseItems{
		ID:              1,
		ItemName:        "Diona",
		ItemDescription: "Cat",
		Type:            "character",
		Stars:           4,
	}
	caseItem3 := data.CaseItems{
		ID:              2,
		ItemName:        "Sword",
		ItemDescription: "Long",
		Type:            "weapon",
		Stars:           3,
	}
	caseItem4 := data.CaseItems{
		ID:              3,
		ItemName:        "Book",
		ItemDescription: "Boring",
		Type:            "weapon",
		Stars:           2,
	}
	caseItem5 := data.CaseItems{
		ID:              4,
		ItemName:        "Stick",
		ItemDescription: "Stinky",
		Type:            "weapon",
		Stars:           1,
	}
	cases = append(cases, caseItem1)
	cases = append(cases, caseItem2)
	cases = append(cases, caseItem3)
	cases = append(cases, caseItem4)
	cases = append(cases, caseItem5)
	fmt.Println(cases)
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/case.gohtml"))
	tmpl.ExecuteTemplate(w, "case.gohtml", cases)
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
	// Accept the metadata struct as a return value.
	items, _, err := app.models.CaseItem.GetAllCaseItem(input.Itemname, input.Type, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/listitems.html"))

	tmpl.ExecuteTemplate(w, "listitems.html", items)
}
