package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/mcfly-svc/mcfly/models"
)

type McflyDB struct {
	*models.DB
}

func NewMcflyDB(databaseUrl string, databaseName string, useSSL bool) *McflyDB {
	modelsDb, err := models.NewDB(databaseUrl, databaseName, useSSL)
	check(err)
	return &McflyDB{modelsDb}
}

func check(err error) {
	switch v := err.(type) {
	case *models.QueryExecError:
		if v != nil {
			panic(v)
		}
	default:
		if err != nil {
			panic(err)
		}
	}
}

func checkDbNotFoundErr(err error, dbName string) {
	pqErr, ok := err.(*pq.Error)
	if ok {
		if pqErr.Code.Name() == "invalid_catalog_name" {
			log.Fatal(
				fmt.Errorf("Database `%s` does not exist. Run `make database` or `mcfly create-db %s`", dbName, dbName),
			)
		}
	}
	check(err)
}

func isDbErr(err *models.QueryExecError, errName string) bool {
	if err == nil {
		return false
	}
	pqErr, ok := err.DbError.(*pq.Error)
	if !ok {
		panic(err)
	}
	if pqErr.Code.Name() == errName {
		return true
	}
	return false
}

func Connect(dbUrl string) *sqlx.DB {
	dbconn, err := sqlx.Connect("postgres", dbUrl)
	check(err)
	return dbconn
}
