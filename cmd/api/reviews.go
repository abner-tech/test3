package main

import (
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
