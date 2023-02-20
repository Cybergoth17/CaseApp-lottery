package dbconnection

import (
	"fmt"
)

func DbConnection() (psqlInformation string) {
	host := "127.0.0.1"
	port := "5432"
	user := "postgres"
	password := "24052004"
	dbname := "final"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return psqlInfo
}
