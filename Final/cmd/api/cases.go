package main

import (
	"Final/internal/data"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) createCaseHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	casename := r.Form.Get("name")
	stars := r.Form.Get("price")
	item1 := r.Form.Get("items1")
	item2 := r.Form.Get("items2")
	item3 := r.Form.Get("items3")
	item4 := r.Form.Get("items4")
	item5 := r.Form.Get("items5")
	caseitems := make([]data.CaseItems, 0)
	x, _ := app.models.CaseItem.GetCaseItemByName(item1)
	x1, _ := app.models.CaseItem.GetCaseItemByName(item2)
	x2, _ := app.models.CaseItem.GetCaseItemByName(item3)
	x3, _ := app.models.CaseItem.GetCaseItemByName(item4)
	x4, _ := app.models.CaseItem.GetCaseItemByName(item5)
	caseitems = append(caseitems, *x)
	caseitems = append(caseitems, *x1)
	caseitems = append(caseitems, *x2)
	caseitems = append(caseitems, *x3)
	caseitems = append(caseitems, *x4)

	price, err := strconv.ParseInt(stars, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	item := &data.Case{
		Name:  casename,
		Price: price,
		Items: caseitems,
	}

	err = app.models.Case.InsertItem(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/items/%d", item.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"case item": item}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCaseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
		return
	}
	item, err := app.models.CaseItem.GetCaseItem(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"case item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCaseHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.CaseItem.DeleteCaseItems(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "item case successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateCaseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Retrieve the movie record as normal.
	item, err := app.models.CaseItem.GetCaseItem(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Use pointers for the Title, Year and Runtime fields.
	var input struct {
		ItemName        *string `json:"itemname"`
		ItemDescription *string `json:"itemdescription"`
		Type            *string `json:"type"`
		Stars           *int64  `json:"stars"`
		Image           *[]byte `json:"image"`
	}
	// Decode the JSON as normal.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.ItemName != nil {
		item.ItemName = *input.ItemName
	}
	// We also do the same for the other fields in the input struct.
	if input.ItemDescription != nil {
		item.ItemDescription = *input.ItemDescription
	}
	if input.Type != nil {
		item.Type = *input.Type
	}
	if input.Stars != nil {
		item.Stars = *input.Stars
	}

	err = app.models.CaseItem.UpdateItem(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"updated case item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
