package main

import (
	"fmt"
	"net/http"
)

func (a *applicationDependences) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()
	a.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (a *applicationDependences) errorResponseJSON(w http.ResponseWriter, r *http.Request, status int, message any) {
	errorData := envelope{"error": message}
	err := a.writeJSON(w, status, errorData, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

// send an error response if our server messes up
func (a *applicationDependences) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//first thing is to log error message
	a.logError(r, err)
	//prepare a response to send to the client
	message := "the server encountered a problem and cound not process your request"
	a.errorResponseJSON(w, r, http.StatusInternalServerError, message)
}

// send an error response of our client messes up with a 404
func (a *applicationDependences) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	//we only log server errors, not client errors
	//prepare a response to send to the client
	message := "the requested resource cound not be found"
	a.errorResponseJSON(w, r, http.StatusNotFound, message)
}

// semd an error response if our client messes up with 405
func (a *applicationDependences) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	//we only log server errors, not client errors
	//prepare a FORMATED response to send to the client
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	a.errorResponseJSON(w, r, http.StatusMethodNotAllowed, message)
}

// send an error response if our client messes up with a 400 (bad request)
func (a *applicationDependences) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
}

func (a *applicationDependences) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	a.errorResponseJSON(w, r, http.StatusUnprocessableEntity, errors)
}

func (a *applicationDependences) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	a.errorResponseJSON(w, r, http.StatusTooManyRequests, message)
}

// send an error response if we have an edit confict status 409
func (a *applicationDependences) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	a.errorResponseJSON(w, r, http.StatusConflict, message)
}

// return 404 unauthorized status code
func (a *applicationDependences) invalidCredentialResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication response"
	a.errorResponseJSON(w, r, http.StatusUnauthorized, message)
}

// We set the WWW-Authenticate header to give a hint to the user as to what they need to provide. Don't want to leave them guessing
func (a *applicationDependences) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("www-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	a.errorResponseJSON(w, r, http.StatusUnauthorized, message)
}
