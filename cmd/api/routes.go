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
	router.HandlerFunc(http.MethodGet, "/api/v1/user/:u_id/lists", a.listUsersReadingLists)

	//the following is the method in which we'll wrap all of our endpoints
	//router.HandlerFunc(http.MethodPost, "/api/v1/SOME ENDPOINT", a.requireActivatedUser(SOME_HANDLER_FUNCTION))

	//READING LISTS SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/lists", a.listAllReadingListHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/lists/:rl_id", a.getSpecificReadingListHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/lists", a.createReadingListHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/lists/:rl_id", a.updateReadingListhandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:rl_id", a.deleteReadingListHander)
	router.HandlerFunc(http.MethodPost, "/api/v1/lists/:rl_id/books", a.addBookToReadingListHandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:rl_id/books", a.deleteBookInReadingListHandler)

	//user section
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:uid", a.listUserProfileHandler)

	//BOOKS SECTION
	router.HandlerFunc(http.MethodGet, "/api/v1/books", a.listAllBooksHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/books", a.addBookHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:b_id", a.listSpecificBookHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/books/:b_id", a.updateBookDetailsHandlers)
	router.HandlerFunc(http.MethodDelete, "/api/v1/books/:b_id", a.deleteBookHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/book/search", a.searchOnBooksHandler)

	//REVIEWS SECTION
	router.HandlerFunc(http.MethodPost, "/api/v1/books/:r_id/reviews", a.addReviewForBooksHandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/reviews/:r_id", a.deleteReviewForBookHandler)
	router.HandlerFunc(http.MethodGet, "/api/v_1/books/:rb_id/reviews", a.listAllReviewsForBookHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/reviews/:r_id", a.updateReviewForBookHandler)

	return a.recoverPanic(a.rateLimiting(a.authenticate(router)))
}
