package data

import (
	"context"
	"database/sql"
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
