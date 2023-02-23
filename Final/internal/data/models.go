package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Case      CaseModel
	Tokens    TokenModel // Add a new Tokens field.
	Users     UserModel
	CaseItem  CaseItemModel
	Inventory InventoryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Case:      CaseModel{DB: db},
		Tokens:    TokenModel{DB: db}, // Initialize a new TokenModel instance.
		Users:     UserModel{DB: db},
		CaseItem:  CaseItemModel{DB: db},
		Inventory: InventoryModel{DB: db},
	}
}
