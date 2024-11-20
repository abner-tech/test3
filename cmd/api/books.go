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

	books, metadata, err := a.bookModel.GetAll(queryParameterData.Fileters)
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

// function to update
func (a *applicationDependences) updateBookDetailsHandlers(w http.ResponseWriter, r *http.Request) {
	//get id parameter
	id, err := a.readIDParam(r, "b_id")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	//fetch previouse record infomration
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

	var incomingData struct {
		Title            *string    `json:"title"`
		Authors          *[]string  `json:"authors"`
		ISBN             *int64     `json:"isbn"`
		Publication_date *time.Time `json:"publication_date"`
		Genre            *[]string  `json:"genre"`
		Description      *string    `json:"description"`
	}

	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	//check for updated fields
	if incomingData.Title != nil {
		book.Title = *incomingData.Title
	}

	if incomingData.Authors != nil {
		book.Authors = *incomingData.Authors
	}
	if incomingData.ISBN != nil {
		book.ISBN = *incomingData.ISBN
	}
	if incomingData.Publication_date != nil {
		book.Publication_Date = *incomingData.Publication_date
	}
	if incomingData.Genre != nil {
		book.Genre = *incomingData.Genre
	}
	if incomingData.Description != nil {
		book.Description = *incomingData.Description
	}

	//validate content sent
	v := validator.New()
	data.ValidateBook(v, book)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.bookModel.UpdateBook(book)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	//display the book
	data := envelope{
		"book": book,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// handler to delete book
func (a *applicationDependences) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	//fetch the provided id
	id, err := a.readIDParam(r, "b_id")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	//delete
	err = a.bookModel.DeleteBook(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	//display the message
	data := envelope{
		"messaeg": "book deleted sucessfully",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// list all books using filters and pagination (optional)
func (a *applicationDependences) searchOnBooksHandler(w http.ResponseWriter, r *http.Request) {
	//to hold query parameters
	var queryParameterData struct {
		Title  string
		Author string
		Genre  string
		data.Fileters
	}

	//get query parameters from url
	queryParameter := r.URL.Query()

	//load the query parameters into the created struct
	queryParameterData.Title = a.getSingleQueryParameter(queryParameter, "title", "")
	queryParameterData.Author = a.getSingleQueryParameter(queryParameter, "author", "")
	queryParameterData.Genre = a.getSingleQueryParameter(queryParameter, "genre", "")
	v := validator.New()

	queryParameterData.Fileters.Page = a.getSingleIntigerParameter(queryParameter, "page", 1, v)
	queryParameterData.Fileters.PageSize = a.getSingleIntigerParameter(queryParameter, "page_size", 10, v)
	queryParameterData.Fileters.Sorting = a.getSingleQueryParameter(queryParameter, "sorting", "id")
	queryParameterData.Fileters.SortSafeList = []string{"id", "title", "author", "genre", "-id", "-title", "-author", "-genre"}

	data.ValidateFilters(v, queryParameterData.Fileters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	books, metadata, err := a.bookModel.SearchGetAll(queryParameterData.Title, queryParameterData.Author, queryParameterData.Genre, queryParameterData.Fileters)
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
