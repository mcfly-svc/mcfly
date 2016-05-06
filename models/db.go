package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Datastore interface {
	SaveProject(*Project) error
	DeleteProject(int64) error
	GetProjects() ([]Project, error)
	GetProjectById(int64) (*Project, error)

	SaveUser(*User) error
	GetUserByAccessToken(string) (*User, error)
	GetUserByProviderToken(*ProviderAccessToken) (*User, error)
	SetUserProviderToken(int64, *ProviderAccessToken) error
	//DeleteUser(int64) error
	//GetUsers() ([]User, error)
	//GetUserById(int64) (*User, error)
}

type DB struct {
	*sqlx.DB
}

func NewDB(dbName string) (*DB, error) {
	db, err := sqlx.Open("postgres", dbName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Get(dest interface{}, query string, args ...interface{}) *QueryExecError {
	err := db.DB.Get(dest, query, args...)
	if err != nil {
		return NewQueryError(query, err, args)
	}
	return nil
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, *QueryExecError) {
	res, err := db.DB.Exec(query, args...)
	if err != nil {
		return res, NewQueryError(query, err, args)
	}
	return res, nil
}

func (db *DB) QueryRowScan(v interface{}, query string, args ...interface{}) *QueryExecError {
	r := db.DB.QueryRow(query, args...)
	if err := r.Scan(v); err != nil {
		return NewQueryError(query, err, args)
	}
	return nil
}
