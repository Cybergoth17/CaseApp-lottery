package main

import (
	"Final/internal/data"
	"fmt"
	"html/template"
	"net/http"
)

var tmpl *template.Template

func (app *application) mainPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/indexx.html"))
	tmpl.ExecuteTemplate(w, "indexx.html", nil)
}
func (app *application) RegPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("cmd/api/templates/reg.html"))
	tmpl.ExecuteTemplate(w, "reg.html", nil)
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
