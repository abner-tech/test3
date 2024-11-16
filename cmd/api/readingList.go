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
		println(err.Error())
		//no record exist for the user id provideds
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
