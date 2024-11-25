package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/abner-tech/Test3-Api.git/internal/validator"
)

// database connection
type ReviewModel struct {
	DB *sql.DB
}

type Review struct {
	ID           int32     `json:"id"`
	Book_ID      int64     `json:"book_id"`
	User_name    string    `json:"user_name"`
	Rating       float32   `json:"rating"`
	ReviewText   string    `json:"review_text"`
	HelpfulCount int32     `json:"helpful_count"`
	Created_at   time.Time `json:"created_at"`
	Version      int16     `json:"version"`
}

func ValidateReview(v *validator.Validator, review *Review) {
	//validate values
	v.Check(review.User_name != "", "user_name", "must be provided")
	v.Check(len(review.User_name) <= 25, "user_name", "must not be more than 25 bytes")

	v.Check(review.Rating >= 0 && review.Rating <= 5, "rating", "must be a number between 1 and 5")

	v.Check(review.ReviewText != "", "review_text", "must be provided")
	v.Check(len(review.ReviewText) <= 100, "review_text", "must not be more than 100 bytes")
}

func (r *ReviewModel) InsertReview(review *Review) error {
	query := `
	INSERT INTO reviews (book_id, user_name, rating, review_text)
	VALUES ($1, $2, $3, $4)
	RETURNING id, helpful_count, created_at, version
	`
	args := []any{review.Book_ID, review.User_name, review.Rating, review.ReviewText}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(
		&review.ID,
		&review.HelpfulCount,
		&review.Created_at,
		&review.Version,
	)
}

func (r *ReviewModel) GetAllReviews(filters Fileters) ([]*Review, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, book_id, user_name, rating, review_text, helpful_count, created_at, version
	FROM reviews
	ORDER BY %s %s, id ASC
	LIMIT $1 OFFSET $2
	`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
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

	reviews := []*Review{}

	for rows.Next() {
		var review Review
		err := rows.Scan(
			&totalRecords,
			&review.ID,
			&review.Book_ID,
			&review.User_name,
			&review.Rating,
			&review.ReviewText,
			&review.HelpfulCount,
			&review.Created_at,
			&review.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		reviews = append(reviews, &review)
	}
	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	//create the metadata
	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return reviews, metadata, nil
}

func (r *ReviewModel) GetByID(id int64) (*Review, error) {
	query := `
	SELECT id, book_id, user_name, rating, review_text, helpful_count, created_at, version
	FROM reviews
	WHERE id = $1
	`

	var review Review

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.Book_ID,
		&review.User_name,
		&review.Rating,
		&review.ReviewText,
		&review.HelpfulCount,
		&review.Created_at,
		&review.Version,
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
	return &review, nil

}

func (r ReviewModel) UpdateReview(review *Review) error {
	query := `
	UPDATE reviews
	SET rating = $1, review_text=$2, version=version+1
	WHERE id=$3 AND book_id=$4 AND version=$5
	RETURNING version
	`
	args := []any{review.Rating, review.ReviewText, review.ID, review.Book_ID, review.Version}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(
		&review.Version,
	)
}
