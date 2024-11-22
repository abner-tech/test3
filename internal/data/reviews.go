package data

import (
	"context"
	"database/sql"
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
