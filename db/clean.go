package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Clean deletes all data from the database
func Clean(dbUrl string) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Cleaning the database")

	tx := db.MustBegin()
	tx.MustExec("DELETE FROM build")
	tx.MustExec("DELETE FROM user_project")
	tx.MustExec("DELETE FROM project")
	tx.MustExec("DELETE FROM provider_access_token")
	tx.MustExec("DELETE FROM marsupi_user")
	tx.Commit()

}
