package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/abner-tech/Test3-Api.git/internal/validator"
	"github.com/lib/pq"
)

// database connection
type BookModel struct {
	DB *sql.DB
}

// type struct
type Book struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	Authors          []string  `json:"author"`
	ISBN             int64     `json:"isbn"`
	Publication_Date time.Time `json:"publication_date"`
	Genre            []string  `json:"genre"`
	Description      string    `json:"description"`
	Average_Rating   float32   `json:"average_rating"`
	Version          int16     `json:"version"`
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 100, "title", "must not be more than 100 bytes long")

	v.Check(len(book.Authors) > 0, "genre", "must have at least one author")
	seenAuthors := make(map[string]bool) // Map to track duplicate authors
	for i, author := range book.Authors {
		// Check length
		v.Check(len(author) <= 25, fmt.Sprintf("authors[%d]", i), "each author name must not be more than 25 bytes long")
		// Check for duplicates
		if seenAuthors[author] {
			v.Check(false, fmt.Sprintf("authors[%d]", i), "author names must not be repeated")
		}
		seenAuthors[author] = true
	}

	v.Check(book.ISBN > 0, "isbn", "must be a positive number")

	v.Check(!book.Publication_Date.IsZero(), "publication_date", "must be provided")
	v.Check(book.Publication_Date.Before(time.Now()), "publication_date", "cannot be in the future")

	v.Check(len(book.Genre) > 0, "genre", "must have at least one genre")
	seenGenre := make(map[string]bool)
	for i, genre := range book.Genre {
		v.Check(len(genre) <= 25, fmt.Sprintf("genre[%d]", i), "each genre must not be more than 50 characters long")

		if seenGenre[genre] {
			v.Check(false, fmt.Sprintf("genre[%d]", i), "genre entries must not be repeated")
		}
		seenAuthors[genre] = true
	}

	v.Check(len(book.Description) > 0, "description", "must be provided")
	v.Check(len(book.Description) <= 500, "description", "must not be more than 500 bytes long")
}

// insert book to db
func (b *BookModel) Insert(book *Book) error {
	query := `
	INSERT INTO books (title, authors, isbn, publication_date, genre, description)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, average_rating
	`

	args := []any{book.Title, pq.Array(book.Authors), book.ISBN, book.Publication_Date, pq.Array(book.Genre), book.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return b.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.Average_Rating)
}

// list all books
func (b *BookModel) GetAll(filters Fileters) ([]*Book, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, title, authors, isbn, publication_date, genre, description, average_rating , version
	FROM books
	ORDER BY %s %s, id ASC
	LIMIT $1 OFFSET $2
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	//check for errors
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
	books := []*Book{}

	for rows.Next() {
		var book Book
		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.Title,
			pq.Array(&book.Authors),
			&book.ISBN,
			&book.Publication_Date,
			pq.Array(&book.Genre),
			&book.Description,
			&book.Average_Rating,
			&book.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		books = append(books, &book)
	}
	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	//create the metadata
	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return books, metadata, nil
}

// fetch from database using id
func (b *BookModel) GetByID(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT title, authors, isbn, publication_date, genre, description, average_rating, version
	FROM books
	WHERE id = $1
	`
	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, id).Scan(
		&book.Title,
		pq.Array(&book.Authors),
		&book.ISBN,
		&book.Publication_Date,
		pq.Array(&book.Genre),
		&book.Description,
		&book.Average_Rating,
		&book.Version,
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
	return &book, nil
}

// update a book record
func (b *BookModel) UpdateBook(book *Book) error {
	query := `
	UPDATE books
	SET title = $1, authors = $2, isbn =$3, publication_date = $4, genre = $5, description = $6
	WHERE id = $7
	RETURNING version
	`
	args := []any{
		book.Title,
		pq.Array(book.Authors),
		book.ISBN,
		book.Publication_Date,
		pq.Array(book.Genre),
		book.Description,
		book.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return b.DB.QueryRowContext(ctx, query, args...).Scan(
		&book.Version,
	)
}

func (b *BookModel) DeleteBook(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM books
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//excecute the query
	result, err := b.DB.ExecContext(ctx, query, id)
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

// list all books
func (b *BookModel) SearchGetAll(title, author, genre string, filters Fileters) ([]*Book, Metadata, error) {

	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, title, authors, isbn, publication_date, genre, description, average_rating, version
	FROM books
	WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
	  AND (to_tsvector('simple', array_to_string(authors, ' ')) @@ plainto_tsquery('simple', $2) OR $2 = '')
	  AND (to_tsvector('simple', array_to_string(genre, ' ')) @@ plainto_tsquery('simple', $3) OR $3 = '')
	ORDER BY %s %s, id ASC
	LIMIT $4 OFFSET $5;
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.DB.QueryContext(ctx, query, title, author, genre, filters.limit(), filters.offset())
	//check for errors
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
	books := []*Book{}

	for rows.Next() {
		var book Book
		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.Title,
			pq.Array(&book.Authors),
			&book.ISBN,
			&book.Publication_Date,
			pq.Array(&book.Genre),
			&book.Description,
			&book.Average_Rating,
			&book.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		books = append(books, &book)
	}
	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	//create the metadata
	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return books, metadata, nil
}

func (b *BookModel) BookExists(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	SELECT id 
	FROM books
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ID int64

	return b.DB.QueryRowContext(ctx, query, id).Scan(&ID)
}
