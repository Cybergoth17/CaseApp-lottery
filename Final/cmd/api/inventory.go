package main

import (
	"Final/internal/data"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) toInventoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	caseitem, _ := app.models.CaseItem.GetCaseItem(num)
	token, _ := r.Cookie("token")
	user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token.Value)
	inventory, _ := app.models.Inventory.GetByID(user.ID)
	inventory.Items = append(inventory.Items, *caseitem)
	update := app.models.Inventory.Update(inventory)
	fmt.Println(update)
	http.Redirect(w, r, "/", 303)
}
