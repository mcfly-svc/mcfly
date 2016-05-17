package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Clean deletes all data from the database
func Clean(conn *sqlx.DB) {
	tx := conn.MustBegin()
	tx.MustExec("DELETE FROM build")
	tx.MustExec("DELETE FROM user_project")
	tx.MustExec("DELETE FROM project")
	tx.MustExec("DELETE FROM provider_access_token")
	tx.MustExec("DELETE FROM marsupi_user")
	tx.Commit()
}
