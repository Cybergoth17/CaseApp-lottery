package main

import (
	"Final/internal/data"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (app *application) createCaseItemsHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	data1, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	itemname := r.Form.Get("itemname")
	itemdesc := r.Form.Get("itemdescription")
	typee := r.Form.Get("type")
	stars := r.Form.Get("stars")
	num, err := strconv.ParseInt(stars, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	item := &data.CaseItems{
		ItemName:        itemname,
		ItemDescription: itemdesc,
		Type:            typee,
		Stars:           num,
		Image:           data1,
	}
	// Initialize a new Validator.

	err = app.models.CaseItem.InsertItem(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// When sending an HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/items/%d", item.ID))

	http.Redirect(w, r, "/list", 303)
}

func (app *application) showItemsHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) deleteItemHandler(w http.ResponseWriter, r *http.Request) {

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
func (app *application) updateItemCaseHandler(w http.ResponseWriter, r *http.Request) {
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
