package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/abner-tech/Test3-Api.git/internal/data"
	"github.com/abner-tech/Test3-Api.git/internal/validator"
)

func (a *applicationDependences) addReviewForBooksHandler(w http.ResponseWriter, r *http.Request) {
	//get id parameter
	book_id, err := a.readIDParam(r, "r_id")
	if err != nil || book_id < 1 {
		a.notFoundResponse(w, r)
		return
	}

	var incomingData struct {
		User_name  string  `json:"user_name"`
		Rating     float32 `json:"rating"`
		ReviewText string  `json:"review_text"`
	}

	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return

	}

	review := &data.Review{
		Book_ID:    book_id,
		User_name:  incomingData.User_name,
		Rating:     incomingData.Rating,
		ReviewText: incomingData.ReviewText,
	}

	//validate
	v := validator.New()
	data.ValidateReview(v, review)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.reviewModel.InsertReview(review)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	//setting location header
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/books/%d/reviews", review.ID))

	//send 201 code and wata
	data := envelope{
		"review": review,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependences) deleteReviewForBookHandler(w http.ResponseWriter, r *http.Request) {

}

func (a *applicationDependences) listAllReviewsForBookHandler(w http.ResponseWriter, r *http.Request) {
	//to hold query parameters
	var queryParameterData struct {
		data.Fileters
	}

	//get query parameters from url
	queryParameter := r.URL.Query()

	v := validator.New()

	queryParameterData.Fileters.Page = a.getSingleIntigerParameter(queryParameter, "page", 1, v)
	queryParameterData.Fileters.PageSize = a.getSingleIntigerParameter(queryParameter, "page_size", 10, v)
	queryParameterData.Fileters.Sorting = a.getSingleQueryParameter(queryParameter, "sorting", "id")
	queryParameterData.Fileters.SortSafeList = []string{"id", "-id"}

	data.ValidateFilters(v, queryParameterData.Fileters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	book_id, err := a.readIDParam(r, "rb_id")
	if err != nil || book_id < 1 {
		a.notFoundResponse(w, r)
		return
	}

	err = a.bookModel.BookExists(book_id)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	reviews, metadata, err := a.reviewModel.GetAllReviews(queryParameterData.Fileters)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
			return
		default:
			a.serverErrorResponse(w, r, err)
			return
		}
	}

	data := envelope{
		"reviews":   reviews,
		"@metadata": metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependences) updateReviewForBookHandler(w http.ResponseWriter, r *http.Request) {
	review_id, err := a.readIDParam(r, "r_id")
	if err != nil || review_id < 1 {
		a.notFoundResponse(w, r)
		return
	}

	review, err := a.reviewModel.GetByID(review_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var incomingData struct {
		Rating     *float32 `json:"rating"`
		ReviewText *string  `json:"review_text"`
	}

	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingData.Rating != nil {
		review.Rating = *incomingData.Rating
	}

	if incomingData.ReviewText != nil {
		review.ReviewText = *incomingData.ReviewText
	}

	v := validator.New()
	data.ValidateReview(v, review)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//continue with update
	err = a.reviewModel.UpdateReview(review)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConfilct):
			a.editConflictResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"review": review,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
