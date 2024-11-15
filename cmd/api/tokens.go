package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/abner-tech/Comments-Api.git/internal/data"
	"github.com/abner-tech/Comments-Api.git/internal/validator"
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
