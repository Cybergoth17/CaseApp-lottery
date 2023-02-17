package main

import (
	"Final/internal/data"
	"Final/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createCaseItemsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ItemName        string `json:"itemname"`
		ItemDescription string `json:"itemdescription"`
		Type            string `json:"type"`
		Stars           int64  `json:"stars"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the input struct to a new Movie struct.
	item := &data.CaseItems{
		ItemName:        input.ItemName,
		ItemDescription: input.ItemDescription,
		Type:            input.Type,
		Stars:           input.Stars,
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

	err = app.writeJSON(w, http.StatusCreated, envelope{"case item": item}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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
	items, metadata, err := app.models.CaseItem.GetAllCaseItem(input.Itemname, input.Type, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Include the metadata in the response envelope.
	err = app.writeJSON(w, http.StatusOK, envelope{"items": items, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
