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

	//REGISTERING USER
	router.HandlerFunc(http.MethodPost, "/api/v1/register/user", a.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/activated", a.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	//the following is the method in which we'll wrap all of our endpoints
	//router.HandlerFunc(http.MethodPost, "/api/v1/SOME ENDPOINT", a.requireActivatedUser(SOME_HANDLER_FUNCTION))

	//READING LISTS SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/lists", a.listAllReadingListHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/lists/:rl_id", a.getSpecificReadingListHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/lists", a.createReadingListHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/lists/:rl_id", a.updateReadingListhandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:rl_id", a.deleteReadingListHander)

	//user section
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:uid", a.listUserProfileHandler)

	//BOOKS SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/books", a.listAllBooksHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/books", a.addBookHandler)

	return a.recoverPanic(a.rateLimiting(a.authenticate(router)))
}
