package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abner-tech/Test3-Api.git/internal/data"
	"github.com/abner-tech/Test3-Api.git/internal/validator"
)

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
