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
	// User activation and authentication
	router.HandlerFunc(http.MethodPut, "/api/v1/users/activated", a.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	// User's reading lists
	router.HandlerFunc(http.MethodGet, "/api/v1/user/:u_id/lists", a.requireActivatedUser(a.listUsersReadingLists))

	// READING LISTS SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/lists", a.requireActivatedUser(a.listAllReadingListHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/lists/:rl_id", a.requireActivatedUser(a.getSpecificReadingListHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/lists", a.requireActivatedUser(a.createReadingListHandler))
	router.HandlerFunc(http.MethodPut, "/api/v1/lists/:rl_id", a.requireActivatedUser(a.updateReadingListhandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:rl_id", a.requireActivatedUser(a.deleteReadingListHander))
	router.HandlerFunc(http.MethodPost, "/api/v1/lists/:rl_id/books", a.requireActivatedUser(a.addBookToReadingListHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:rl_id/books", a.requireActivatedUser(a.deleteBookInReadingListHandler))

	// USER SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:uid", a.requireActivatedUser(a.listUserProfileHandler))

	// BOOKS SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/books", a.requireActivatedUser(a.listAllBooksHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/books", a.requireActivatedUser(a.addBookHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:b_id", a.requireActivatedUser(a.listSpecificBookHandler))
	router.HandlerFunc(http.MethodPut, "/api/v1/books/:b_id", a.requireActivatedUser(a.updateBookDetailsHandlers))
	router.HandlerFunc(http.MethodDelete, "/api/v1/books/:b_id", a.requireActivatedUser(a.deleteBookHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/book/search", a.requireActivatedUser(a.searchOnBooksHandler))

	// REVIEWS SECTION
	router.HandlerFunc(http.MethodPost, "/api/v1/books/:r_id/reviews", a.requireActivatedUser(a.addReviewForBooksHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/reviews/:r_id", a.requireActivatedUser(a.deleteReviewForBookHandler))
	router.HandlerFunc(http.MethodGet, "/api/v_1/books/:rb_id/reviews", a.requireActivatedUser(a.listAllReviewsForBookHandler))
	router.HandlerFunc(http.MethodPut, "/api/v1/reviews/:r_id", a.requireActivatedUser(a.updateReviewForBookHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/user/:u_id/reviews", a.requireActivatedUser(a.fetchReviewByIdHandler))

	return a.recoverPanic(a.rateLimiting(a.authenticate(router)))
}
