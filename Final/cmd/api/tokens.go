package main

import (
	"Final/internal/data"
	"Final/internal/validator"
	"errors"
	"net/http"
	"time"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	password := r.Form.Get("password")
	email := r.Form.Get("email")
	// Validate the email and password provided by the client.
	v := validator.New()
	data.ValidateEmail(v, email)
	data.ValidatePasswordPlaintext(v, password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Lookup the user record based on the email address. If no matching user was
	// found, then we call the app.invalidCredentialsResponse() helper to send a 401
	// Unauthorized response to the client (we will create this helper in a moment).
	user, err := app.models.Users.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// If the passwords don't match, then we call the app.invalidCredentialsResponse()
	// helper again and return.
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	cookie := http.Cookie{
		Name:    "token",
		Value:   token.Plaintext,
		Expires: token.Expiry,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", 303)
}
