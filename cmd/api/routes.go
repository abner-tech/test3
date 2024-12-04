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
	router.HandlerFunc(http.MethodGet, "/api/v1/user/:u_id/lists", a.requireActivatedUser(a.requirePermission("reading_list:read", a.listUsersReadingLists)))

	// READING LISTS SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/lists", a.requireActivatedUser(a.requirePermission("reading_list:read", a.listAllReadingListHandler)))
	router.HandlerFunc(http.MethodGet, "/api/v1/lists/:rl_id", a.requireActivatedUser(a.requirePermission("reading_list:read", a.getSpecificReadingListHandler)))
	router.HandlerFunc(http.MethodPost, "/api/v1/lists", a.requireActivatedUser(a.requirePermission("reading_list:write", a.createReadingListHandler)))
	router.HandlerFunc(http.MethodPut, "/api/v1/lists/:rl_id", a.requireActivatedUser(a.requirePermission("reading_list:write", a.updateReadingListhandler)))
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:rl_id", a.requireActivatedUser(a.requirePermission("reading_list:write", a.deleteReadingListHander)))
	router.HandlerFunc(http.MethodPost, "/api/v1/lists/:rl_id/books", a.requireActivatedUser(a.requirePermission("reading_list:write", a.addBookToReadingListHandler)))
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:rl_id/books", a.requireActivatedUser(a.requirePermission("reading_list:write", a.deleteBookInReadingListHandler)))

	// USER SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:uid", a.requireActivatedUser(a.requirePermission("users:read", a.listUserProfileHandler)))

	// BOOKS SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/books", a.requireActivatedUser(a.requirePermission("books:read", a.listAllBooksHandler)))
	router.HandlerFunc(http.MethodPost, "/api/v1/books", a.requireActivatedUser(a.requirePermission("books:write", a.addBookHandler)))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:b_id", a.requireActivatedUser(a.requirePermission("books:read", a.listSpecificBookHandler)))
	router.HandlerFunc(http.MethodPut, "/api/v1/books/:b_id", a.requireActivatedUser(a.requirePermission("books:write", a.updateBookDetailsHandlers)))
	router.HandlerFunc(http.MethodDelete, "/api/v1/books/:b_id", a.requireActivatedUser(a.requirePermission("books:write", a.deleteBookHandler)))
	router.HandlerFunc(http.MethodGet, "/api/v1/book/search", a.requireActivatedUser(a.requirePermission("books:read", a.searchOnBooksHandler)))

	// REVIEWS SECTION
	router.HandlerFunc(http.MethodPost, "/api/v1/books/:r_id/reviews", a.requireActivatedUser(a.requirePermission("reviews:write", a.addReviewForBooksHandler)))
	router.HandlerFunc(http.MethodDelete, "/api/v1/reviews/:r_id", a.requireActivatedUser(a.requirePermission("reviews:write", a.deleteReviewForBookHandler)))
	router.HandlerFunc(http.MethodGet, "/api/v_1/books/:rb_id/reviews", a.requireActivatedUser(a.requirePermission("reviews:read", a.listAllReviewsForBookHandler)))
	router.HandlerFunc(http.MethodPut, "/api/v1/reviews/:r_id", a.requireActivatedUser(a.requirePermission("reviews:write", a.updateReviewForBookHandler)))
	router.HandlerFunc(http.MethodGet, "/api/v1/user/:u_id/reviews", a.requireActivatedUser(a.requirePermission("reviews:read", a.fetchReviewByIdHandler)))

	return a.recoverPanic(a.rateLimiting(a.authenticate(router)))
}
