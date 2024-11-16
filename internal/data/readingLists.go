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
