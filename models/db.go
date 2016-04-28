package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
)

type Datastore interface {
	SaveProject(*Project) error
	DeleteProject(int64) error
	GetProjects() ([]Project, error)
	GetProjectById(int64) (*Project, error)

	SaveUser(*User) error
	GetUserByProviderToken(ProviderAccessToken) (*User, error)
	SetUserProviderToken(int64, ProviderAccessToken) error
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

func handleQueryError(method string, query string, queryError error) *QueryExecError {
	err, ok := queryError.(*pq.Error)
	if !ok {
		log.Fatal(fmt.Sprintf("handleQueryError failed: err `%s` is not a *pq.Error", err.Error()))
	}
	return &QueryExecError{method, query, err, err.Code.Name()}
}
