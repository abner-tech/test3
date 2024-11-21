package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/abner-tech/Test3-Api.git/internal/data"
	"github.com/abner-tech/Test3-Api.git/internal/validator"
)

// create a reading list for the user
func (a *applicationDependences) createReadingListHandler(w http.ResponseWriter, r *http.Request) {
	//create a struct to hold a list
	var incomingData struct {
		ListName        string `json:"name"`
		ListDescription string `json:"description"`
		CreatedBy       int64  `json:"created_by"`
	}

	//decoding
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	reading_List := &data.Reading_List{
		Name:        incomingData.ListName,
		Description: incomingData.ListDescription,
		CreatedBy:   incomingData.CreatedBy,
	}

	//validate inserted data
	v := validator.New()
	data.ValidateReadingList(v, reading_List)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//check if user the list is created for does exists
	err = a.userModel.UserExist(incomingData.CreatedBy)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	//create the list in the database
	err = a.readingListModel.CreateReadingList(reading_List)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	//setting location header path to newly created list
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/lists/%d", reading_List.ID))

	//send json response woth a 201 status code (new resource ccreated)
	data := envelope{
		"readingLists": reading_List,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// fetches all existing reading list for all users
func (a *applicationDependences) listAllReadingListHandler(w http.ResponseWriter, r *http.Request) {
	//to hold query parameters
	var queryParameterData struct {
		Description string
		data.Fileters
	}

	//get query parameters from url
	queryParameter := r.URL.Query()

	//load the query parameters into the created struct
	queryParameterData.Description = a.getSingleQueryParameter(queryParameter, "description", "")
	v := validator.New()

	queryParameterData.Fileters.Page = a.getSingleIntigerParameter(queryParameter, "page", 1, v)
	queryParameterData.Fileters.PageSize = a.getSingleIntigerParameter(queryParameter, "page_size", 10, v)
	queryParameterData.Fileters.Sorting = a.getSingleQueryParameter(queryParameter, "sorting", "id")
	queryParameterData.Fileters.SortSafeList = []string{"id", "created_at", "-id", "-created_at"}

	//check validity of filters
	data.ValidateFilters(v, queryParameterData.Fileters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	readingList, metadata, err := a.readingListModel.GetAll(queryParameterData.Description, queryParameterData.Fileters)
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
		"readingList": readingList,
		"@metadata":   metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// fetches specific reading list using id
func (a *applicationDependences) getSpecificReadingListHandler(w http.ResponseWriter, r *http.Request) {
	//get id parameter
	id, err := a.readIDParam(r, "rl_id")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	list, err := a.readingListModel.GetByID(id)
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
		"reading list": list,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// update a specific reading list
func (a *applicationDependences) updateReadingListhandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r, "rl_id")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	list, err := a.readingListModel.GetByID(id)
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
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	//cheking which field has been updated
	if incomingData.Name != nil {
		list.Name = *incomingData.Name
	}
	if incomingData.Description != nil {
		list.Description = *incomingData.Description
	}

	//validate new list values
	v := validator.New()
	data.ValidateReadingList(v, list)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//proceed with updating record
	err = a.readingListModel.UpdateReadingList(list)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"reading_list": list,
	}

	a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// delete a specific reading list
func (a *applicationDependences) deleteReadingListHander(w http.ResponseWriter, r *http.Request) {
	//fetch the provided id
	id, err := a.readIDParam(r, "rl_id")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	//delete
	err = a.readingListModel.DeleteSingleList(id)
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
		"messaeg": "list deleted sucessfully",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// add a book to users reading list
func (a *applicationDependences) addBookToReadingListHandler(w http.ResponseWriter, r *http.Request) {
	//create a struct to hold a list
	var incomingData struct {
		Book_ID int64  `json:"book_id"`
		Status  string `json:"status"`
	}

	//get list id parameter
	id, err := a.readIDParam(r, "rl_id")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	//decode
	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	bookInList := &data.BookInList{
		Reading_List_ID: id,
		Book_ID:         incomingData.Book_ID,
		Status:          incomingData.Status,
	}

	//validate status
	v := validator.New()
	data.ValidateReadingStatus(v, incomingData.Status)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//check if reading list exist
	err = a.readingListModel.ReadingListExist(id)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	//check if book exist
	err = a.bookModel.BookExists(incomingData.Book_ID)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	//procede to insert to DB
	err = a.readingListModel.AddBookToReadingList(bookInList)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateBookInList):
			v.AddError("book", data.ErrDuplicateBookInList.Error())
			a.failedValidationResponse(w, r, v.Errors)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/lists/%d/books", incomingData.Book_ID))

	data := envelope{
		"added_Book": bookInList,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}
