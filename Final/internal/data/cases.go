package data

import (
	"database/sql"
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
