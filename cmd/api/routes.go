package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependences) routes() http.Handler {
	//setup a new router
	router := httprouter.New()

	//handle 405
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	//method 404
	router.NotFound = http.HandlerFunc(a.notFoundResponse)

	//setup routes
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", a.healthChechHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/register/user", a.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/activated", a.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	//the following is the method in which we'll wrap all of our endpoints
	//router.HandlerFunc(http.MethodPost, "/api/v1/SOME ENDPOINT", a.requireActivatedUser(SOME_HANDLER_FUNCTION))

	return a.recoverPanic(a.rateLimiting(a.authenticate(router)))
}
