package models


import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type Datastore interface {
	SaveProject(*Project) (error)
	DeleteProject(int64) (error)
	GetProjects() ([]Project, error)
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
