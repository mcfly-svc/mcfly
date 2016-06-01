package db

import _ "github.com/lib/pq"

// Clean deletes all data from the database
func (mdb *McflyDB) Clean() {
	tx, err := mdb.Begin()
	checkDbNotFoundErr(err, mdb.DatabaseName)
	tx.Exec("DELETE FROM build")
	tx.Exec("DELETE FROM user_project")
	tx.Exec("DELETE FROM project")
	tx.Exec("DELETE FROM provider_access_token")
	tx.Exec("DELETE FROM mcfly_user")
	tx.Commit()
}
