package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Case struct {
	ID    int64       `json:"id"`
	Name  string      `json:"name"`
	Price int64       `json:"price"`
	Items []CaseItems `json:"items"`
}

type CaseModel struct {
	DB *sql.DB
}

func (m CaseModel) InsertItem(item *Case) error {

	itemsJSON, err := json.Marshal(item.Items)
	if err != nil {
		return err
	}

	query := `
INSERT INTO cases(name, price, items)
VALUES ($1, $2, $3::jsonb)
RETURNING id`
	args := []any{item.Name, item.Price, itemsJSON}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&item.ID)
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

func (m CaseModel) GetCaseID(id int64) (*Case, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT id, name, price, items 
FROM cases 
WHERE id = $1`
	var itemsJSON string
	var item Case
	err := m.DB.QueryRow(query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Price,
		&itemsJSON,
	)
	err = json.Unmarshal([]byte(itemsJSON), &item.Items)
	if err != nil {
		return nil, err
	}

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

func (m CaseModel) GetAllCase(itemname string, typee string, filters Filters) ([]*Case, Metadata, error) {
	// Update the SQL query to include the LIMIT and OFFSET clauses with placeholder
	// parameter values.
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, name,price,items
FROM cases
WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') 
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
	var itemsJSON string
	totalRecords := 0
	items := []*Case{}
	for rows.Next() {
		var item Case
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&item.ID,
			&item.Name,
			&item.Price,
			&itemsJSON,
		)
		err = json.Unmarshal([]byte(itemsJSON), &item.Items)

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
