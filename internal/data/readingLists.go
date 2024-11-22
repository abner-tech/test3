package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/abner-tech/Test3-Api.git/internal/validator"
)

type ReadingListModel struct {
	DB *sql.DB
}

// reading_list type declaration
type Reading_List struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"-"`
	Version     int16     `json:"version"`
}

type BookInList struct {
	Reading_List_ID int64     `json:"reading_list_id"`
	Book_ID         int64     `json:"book_id"`
	Status          string    `json:"status"`
	Created_at      time.Time `json:"created_at"`
	Version         int16     `json:"version"`
}

// validate provided content for list being created
func ValidateReadingList(v *validator.Validator, reading_List *Reading_List) {
	//check if name is empty or too long
	v.Check(reading_List.Name != "", "name", "must be provided")
	v.Check(len(reading_List.Name) <= 25, "name", "must not be more than 25 bytes long")

	//verifying
	v.Check(reading_List.Description != "", "description", "must not be empty")
	v.Check(len(reading_List.Description) <= 250, "description", "must not be more than 250 bytes long")
}

// validate if status for book being added to reading list is correct
func ValidateReadingStatus(v *validator.Validator, readingStatus string) {
	v.Check(readingStatus != "", "status", "must be provided")
	v.Check(readingStatus == "currently reading" || readingStatus == "completed",
		"status",
		"status must be of values 'completed' or 'currently reading'")
}

// create the list for the user
func (r *ReadingListModel) CreateReadingList(reading_List *Reading_List) error {
	query := `
	INSERT INTO reading_lists(name, description,created_by)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version
	`
	args := []any{reading_List.Name, reading_List.Description, reading_List.CreatedBy}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(
		&reading_List.ID,
		&reading_List.CreatedAt,
		&reading_List.Version,
	)
}

// fetch all reading lists for all users PAGINATION used
func (r *ReadingListModel) GetAll(description string, filters Fileters) ([]*Reading_List, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, name, description, created_at, created_by, version
	FROM reading_lists
	WHERE (to_tsvector('simple',description) @@
		plainto_tsquery('simple', $1) OR $1 = '')
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3
	`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, description, filters.limit(), filters.offset())
	//checking for errors
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, Metadata{}, err
		default:
			return nil, Metadata{}, err
		}
	}

	defer rows.Close()

	totalRecords := 0

	lists := []*Reading_List{}

	for rows.Next() {
		var rec Reading_List
		err := rows.Scan(
			&totalRecords,
			&rec.ID,
			&rec.Name,
			&rec.Description,
			&rec.CreatedAt,
			&rec.CreatedBy,
			&rec.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		lists = append(lists, &rec)
	}
	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	//create the metadata
	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)

	return lists, metadata, nil
}

// fetch at most 1 reading list using id
func (r *ReadingListModel) GetByID(id int64) (*Reading_List, error) {
	//check for valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//query
	query := `
	SELECT id, name, description, created_by, created_at, version
	FROM reading_lists
	WHERE id = $1
	`

	//variable to hold row
	var list Reading_List

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&list.ID,
		&list.Name,
		&list.Description,
		&list.CreatedBy,
		&list.CreatedAt,
		&list.Version,
	)

	//check if errors
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &list, nil
}

// update a reading list using list id
func (r *ReadingListModel) UpdateReadingList(reading_List *Reading_List) error {
	query := `
	UPDATE reading_lists
	SET name = $1, description = $2, version = version + 1
	WHERE id = $3
	RETURNING version
	`

	args := []any{reading_List.Name, reading_List.Description, reading_List.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(
		&reading_List.Version,
	)
}

// deletre a reading list from the
func (r *ReadingListModel) DeleteSingleList(id int64) error {
	//validate id
	if id < 1 {
		return ErrRecordNotFound
	}

	//query
	query := `
	DELETE FROM reading_lists
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//excecute the query
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	//check if any rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound //no rows affected
	}
	return nil
}

// check if reading list exists
func (b *ReadingListModel) ReadingListExist(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	SELECT id 
	FROM reading_lists
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ID int64

	return b.DB.QueryRowContext(ctx, query, id).Scan(&ID)
}

// adding book to reading list
func (b *ReadingListModel) AddBookToReadingList(book *BookInList) error {
	query := `
	INSERT INTO reading_list_books(reading_list_id, book_id, status)
	VALUES($1, $2, $3)
	RETURNING created_at, version
	`
	args := []any{book.Reading_List_ID, book.Book_ID, book.Status}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(
		&book.Created_at,
		&book.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "reading_list_books_pkey"`:
			return ErrDuplicateBookInList
		default:
			return err
		}
	}

	return nil
}

func (b *ReadingListModel) DeleteBookFromReadingList(bookID, listID int64) error {
	query := `
	DELETE FROM reading_list_books
	WHERE reading_list_id = $1 AND book_id = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//excecute the query
	result, err := b.DB.ExecContext(ctx, query, listID, bookID)
	if err != nil {
		return err
	}

	//check if any rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound //no rows affected
	}
	return nil

}
