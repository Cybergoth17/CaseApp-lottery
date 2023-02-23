package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type CaseItems struct {
	ID              int64  `json:"id"`
	ItemName        string `json:"itemname"`
	ItemDescription string `json:"itemdescription"`
	Type            string `json:"type"`
	Stars           int64  `json:"stars"`
	Image           []byte `json:"image"`
}

type CaseItemModel struct {
	DB *sql.DB
}

// that we did when creating a movie.
func (m CaseItemModel) InsertItem(item *CaseItems) error {
	query := `
INSERT INTO caseitems(itemname, itemdesc, type, stars,image)
VALUES ($1, $2, $3, $4,$5)
RETURNING id`
	args := []any{item.ItemName, item.ItemDescription, item.Type, item.Stars, item.Image}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&item.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m CaseItemModel) GetCaseItem(id int64) (*CaseItems, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT id, itemname, itemdesc, type,stars,image
FROM caseitems
WHERE id = $1`

	var item CaseItems
	err := m.DB.QueryRow(query, id).Scan(
		&item.ID,
		&item.ItemName,
		&item.ItemDescription,
		&item.Type,
		&item.Stars,
		&item.Image,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &item, nil
}

func (m CaseItemModel) DeleteCaseItems(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
DELETE FROM caseitems
WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m CaseItemModel) UpdateItem(item *CaseItems) error {
	query := `
UPDATE caseitems
SET itemname = $1, itemdesc = $2, type = $3, stars = $4, image=$6
WHERE id = $5
RETURNING id`

	args := []interface{}{
		item.ItemName,
		item.ItemDescription,
		item.Type,
		item.Stars,
		item.ID,
		item.Image,
	}

	return m.DB.QueryRow(query, args...).Scan(&item.ID)
}

func (m CaseItemModel) GetAllCaseItem(itemname string, typee string, filters Filters) ([]*CaseItems, Metadata, error) {
	// Update the SQL query to include the LIMIT and OFFSET clauses with placeholder
	// parameter values.
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, itemname,itemdesc,type,stars, image
FROM caseitems
WHERE (to_tsvector('simple', itemname) @@ plainto_tsquery('simple', $1) OR $1 = '') 
ORDER BY %s %s, id ASC
LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{itemname, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	items := []*CaseItems{}
	for rows.Next() {
		var item CaseItems
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&item.ID,
			&item.ItemName,
			&item.ItemDescription,
			&item.Type,
			&item.Stars,
			&item.Image,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		items = append(items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return items, metadata, nil
}

func (m CaseItemModel) GetCaseItemByName(name string) (*CaseItems, error) {

	if name == "" {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT id, itemname, itemdesc, type,stars,image
FROM caseitems
WHERE itemname = $1`

	var item CaseItems
	err := m.DB.QueryRow(query, name).Scan(
		&item.ID,
		&item.ItemName,
		&item.ItemDescription,
		&item.Type,
		&item.Stars,
		&item.Image,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &item, nil
}
