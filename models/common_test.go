package models_test

import (
	"testing"

	"github.com/mikec/msplapi/db"
	"github.com/mikec/msplapi/models"
)

var (
	dbUrl string
	DB    *models.DB
)

func getDB() *models.DB {
	return DB
}

func init() {
	dbUrl = "postgres://localhost:5432/marsupi_test?sslmode=disable"
	newDb, err := models.NewDB(dbUrl)
	if err != nil {
		panic(err)
	}
	DB = newDb
}

func resetDB() {
	cleanupDB()
	seedDB()
}

func cleanupDB() {
	db.Clean(dbUrl)
}

func seedDB() {
	db.Seed(dbUrl)
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
