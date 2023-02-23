package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Inventory struct {
	ID     int64       `json:"id"`
	UserID int64       `json:"-"`
	Items  []CaseItems `json:"items"`
}

func generateInventory(userID int64) (*Inventory, error) {
	items := make([]CaseItems, 0)

	x := CaseItems{
		ID:              0,
		ItemName:        "Suprise",
		ItemDescription: "",
		Type:            "",
		Stars:           0,
		Image:           nil,
	}
	items = append(items, x)
	inventory := &Inventory{
		ID:     userID,
		UserID: userID,
		Items:  items,
	}
	return inventory, nil
}

type InventoryModel struct {
	DB *sql.DB
}

func (m InventoryModel) New(userID int64) error {
	inventory, err := generateInventory(userID)
	if err != nil {
		return err
	}
	err = m.Insert(inventory)
	return err
}

func (m InventoryModel) Insert(inventory *Inventory) error {
	itemsJSON, err := json.Marshal(inventory.Items)
	if err != nil {
		return err
	}
	query := `
INSERT INTO inventory ( ID,user_id, items)
VALUES ($1,$2, $3::jsonb )`
	args := []interface{}{inventory.ID, inventory.UserID, itemsJSON}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m InventoryModel) Update(inventory *Inventory) error {
	itemsJSON, err := json.Marshal(inventory.Items)
	if err != nil {
		return err
	}
	query := `
UPDATE inventory
SET ID = $1, user_id = $2, items = $3::jsonb
WHERE ID = $1`
	args := []interface{}{inventory.ID, inventory.UserID, itemsJSON}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = m.DB.ExecContext(ctx, query, args...)
	return err
}
func (m InventoryModel) Check(userID int64) bool {
	var hasInventory bool
	err := m.DB.QueryRow("SELECT COUNT(*) > 0 FROM inventory WHERE user_id = $1", userID).Scan(&hasInventory)
	if err != nil {
		fmt.Println(err)
	}
	if hasInventory {
		return true
	} else {
		return false
	}
}

func (m InventoryModel) GetByID(id int64) (*Inventory, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT id, user_id, items
FROM inventory
WHERE id = $1`
	var itemsJSON string
	var item Inventory
	err := m.DB.QueryRow(query, id).Scan(
		&item.ID,
		&item.UserID,
		&itemsJSON,
	)
	err = json.Unmarshal([]byte(itemsJSON), &item.Items)
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
