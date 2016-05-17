package models_test

import (
	"testing"

	"github.com/mikec/msplapi/db"
	"github.com/mikec/msplapi/models"
)

var (
	DB *models.DB
)

func getDB() *models.DB {
	return DB
}

func init() {
	newDb, err := models.NewDB("postgres://localhost:5432/marsupi_test?sslmode=disable")
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
	db.Clean(DB.DB)
}

func seedDB() {
	db.Seed(DB.DB)
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
