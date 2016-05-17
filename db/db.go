package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func isDbErr(err error, errName string) bool {
	if err == nil {
		return false
	}
	pqErr, ok := err.(*pq.Error)
	if !ok {
		log.Fatalln(err)
	}
	if pqErr.Code.Name() == errName {
		return true
	}
	return false
}

func Connect(dbUrl string) *sqlx.DB {
	dbconn, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}
	return dbconn
}
