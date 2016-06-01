package db

import "fmt"

// Creates a new database
func (mdb *McflyDB) Create(databaseName string) {
	_, err := mdb.Exec(fmt.Sprintf("CREATE DATABASE %s", databaseName))
	if err != nil {
		if isDbErr(err, "duplicate_database") {
			fmt.Printf("Database %s already exists!\n", databaseName)
		} else {
			check(err)
		}
	}
}
