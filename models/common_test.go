package models_test

import (
	"testing"

	"github.com/mikec/msplapi/db"
)

var (
	mdb *db.MsplDB
)

func getMsplDB() *db.MsplDB {
	return mdb
}

func init() {
	cfg := GetTestConfig()
	mdb = db.NewMsplDB(cfg.DatabaseUrl, cfg.DatabaseName, cfg.DatabaseUseSSL)
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
