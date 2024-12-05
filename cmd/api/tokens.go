package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/abner-tech/Test3-Api.git/internal/data"
	"github.com/abner-tech/Test3-Api.git/internal/validator"
)

func (a *applicationDependences) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, incomingData.Email)
	data.ValidatePassword(v, incomingData.Password)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//checking if user exist for email
	user, err := a.userModel.GetByEmail(incomingData.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.invalidCredentialResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	//user is found so we verify password
	match, err := user.Password.Matches(incomingData.Password)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	//wrong password
	if !match {
		a.invalidCredentialResponse(w, r)
		return
	}

	token, err := a.tokenModel.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"authentication_token": token,
	}

	err = a.writeJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return //just in case to terminate function
	}

}

func (a *applicationDependences) createPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Email string `json:"email"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, incomingData.Email)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//checking if user exist for email
	user, err := a.userModel.GetByEmail(incomingData.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.invalidCredentialResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := a.tokenModel.New(user.ID, 30*time.Minute, data.ScopePasswordReset)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"message": "temporary token has been sent to your email",
	}

	a.background(func() {
		data := map[string]any{
			"passwordResetToken": token.PlainText,
			"userID":             user.ID,
			"tokenExpiryTime":    token.Expiry,
		}
		err = a.mailer.Send(user.Email, "user_reset_password.tmpl", data)
		if err != nil {
			a.logger.Error(err.Error())
		}
	})

	err = a.writeJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return //just in case to terminate function
	}

}
