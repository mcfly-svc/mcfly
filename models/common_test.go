package models_test

import (
	"testing"

	"github.com/mcfly-svc/mcfly/db"
)

var (
	mdb *db.McflyDB
)

func getMcflyDB() *db.McflyDB {
	return mdb
}

func init() {
	cfg := GetTestConfig()
	mdb = db.NewMcflyDB(cfg.DatabaseUrl, cfg.DatabaseName, cfg.DatabaseUseSSL)
}

func resetDB() {
	cleanupDB()
	seedDB()
}

func cleanupDB() {
	mdb.Clean()
}

func seedDB() {
	mdb.Seed()
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
