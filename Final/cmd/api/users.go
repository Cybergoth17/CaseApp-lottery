package main

import (
	"Final/internal/data"
	"Final/internal/validator"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	name := r.Form.Get("username")
	password := r.Form.Get("password")
	email := r.Form.Get("email")
	fmt.Println(password)
	fmt.Println(email)

	user := &data.User{
		Name:      name,
		Email:     email,
		Activated: false,
		Role:      "admin",
		Balance:   5000,
	}

	err = user.Password.Set(password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		if v.Errors["password"] != "" {
			http.Redirect(w, r, "/errors?err="+v.Errors["password"], 303)
			return
		} else if v.Errors["name"] != "" {
			http.Redirect(w, r, "/errors?err="+v.Errors["name"], 303)
			return
		}
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {

		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.background(func() {

		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	http.Redirect(w, r, "/confirm", 303)
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	token := app.readString(qs, "token", "")

	v := validator.New()
	if data.ValidateTokenPlaintext(v, token); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true
	user.Balance = 5000
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/", 303)
}

func (app *application) balanceMinusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	x, _ := app.models.Case.GetCaseID(num)
	token, _ := r.Cookie("token")
	userFound, _ := app.models.Users.GetForToken(data.ScopeAuthentication, token.Value)
	balanceNum := userFound.Balance
	if balanceNum <= 0 || balanceNum-x.Price < 0 {
		//todo
		http.Redirect(w, r, "/", 303)
	}
	newB := balanceNum - x.Price

	userFound.Balance = newB
	if newB > 0 {
		_ = app.models.Users.Update(userFound)
		http.Redirect(w, r, "/case/"+id, 303)
	}

}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = app.models.Users.DeleteUser(num)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	token, _ := r.Cookie("token")
	token.MaxAge = -1
	http.SetCookie(w, token)
	http.Redirect(w, r, "/logout", 303)
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	token, _ := r.Cookie("token")
	user, _ := app.models.Users.GetForToken(data.ScopeAuthentication, token.Value)
	user.Name = r.Form.Get("name")
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	http.Redirect(w, r, "/profile", 303)
}
