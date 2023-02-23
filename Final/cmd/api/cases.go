package main

import (
	"Final/internal/data"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) createCaseHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	caseitems := make([]data.CaseItems, 0)
	for i := 1; i <= 5; i++ {
		itemName := r.FormValue(fmt.Sprintf("star%d", i))
		if itemName != "" {
			item, err := app.models.CaseItem.GetCaseItemByName(itemName)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			caseitems = append(caseitems, *item)
		}
	}

	// Parse the price from the form data
	price, err := strconv.ParseInt(r.Form.Get("price"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	// Create the case item
	item := &data.Case{
		Name:  r.Form.Get("name"),
		Price: price,
		Items: caseitems,
	}

	// Insert the case item into the database
	err = app.models.Case.InsertItem(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Redirect to the list of case items
	http.Redirect(w, r, "/list", 303)
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

	vars := mux.Vars(r)
	id := vars["id"]
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = app.models.Case.DeleteCase(num)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	http.Redirect(w, r, "/case", 303)
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
