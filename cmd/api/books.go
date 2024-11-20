package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/abner-tech/Test3-Api.git/internal/data"
	"github.com/abner-tech/Test3-Api.git/internal/validator"
)

// insert book to db
func (a *applicationDependences) addBookHandler(w http.ResponseWriter, r *http.Request) {

	//var to hold user input data
	var incomingData struct {
		Title            string    `json:"title"`
		Authors          []string  `json:"author"`
		ISBN             int64     `json:"isbn"`
		Publication_Date time.Time `json:"publication_date"`
		Genre            []string  `json:"genre"`
		Description      string    `json:"description"`
	}

	//parse info
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
	}

	book := &data.Book{
		Title:            incomingData.Title,
		Authors:          incomingData.Authors,
		ISBN:             incomingData.ISBN,
		Publication_Date: incomingData.Publication_Date,
		Genre:            incomingData.Genre,
		Description:      incomingData.Description,
	}

	//validate content sent
	v := validator.New()
	data.ValidateBook(v, book)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.bookModel.Insert(book)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	//setting location header
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/books/%d", book.ID))

	//send 201 code and wata
	data := envelope{
		"book": book,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// list all books using filters and pagination (optional)
func (a *applicationDependences) listAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	//to hold query parameters
	var queryParameterData struct {
		Title       string
		Description string
		data.Fileters
	}

	//get query parameters from url
	queryParameter := r.URL.Query()

	//load the query parameters into the created struct
	queryParameterData.Title = a.getSingleQueryParameter(queryParameter, "title", "")
	queryParameterData.Description = a.getSingleQueryParameter(queryParameter, "description", "")
	v := validator.New()

	queryParameterData.Fileters.Page = a.getSingleIntigerParameter(queryParameter, "page", 1, v)
	queryParameterData.Fileters.PageSize = a.getSingleIntigerParameter(queryParameter, "page_size", 10, v)
	queryParameterData.Fileters.Sorting = a.getSingleQueryParameter(queryParameter, "sorting", "id")
	queryParameterData.Fileters.SortSafeList = []string{"id", "author", "-id", "-author"}

	data.ValidateFilters(v, queryParameterData.Fileters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	books, metadata, err := a.bookModel.GetAll(queryParameterData.Title, queryParameterData.Description, queryParameterData.Fileters)
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
		"books":     books,
		"@metadata": metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// list 1 book using is
func (a *applicationDependences) listSpecificBookHandler(w http.ResponseWriter, r *http.Request) {
	//get id parameter
	id, err := a.readIDParam(r, "b_id")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	book, err := a.bookModel.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	//display the list
	data := envelope{
		"book": book,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}
