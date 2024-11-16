package data

import (
	"context"
	"database/sql"
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

// validate provided content for list being created
func ValidateReadingList(v *validator.Validator, reading_List *Reading_List) {
	//check if name is empty or too long
	v.Check(reading_List.Name != "", "name", "must be provided")
	v.Check(len(reading_List.Name) <= 25, "name", "must not be more than 25 bytes long")

	//verifying
	v.Check(reading_List.Description != "", "description", "must not be empty")
	v.Check(len(reading_List.Description) <= 250, "description", "must not be more than 250 bytes long")
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
