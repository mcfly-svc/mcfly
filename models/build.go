package models

import "fmt"

type Build struct {
	ID        int64  `db:"id" json:"id"`
	ProjectID int64  `db:"provider_id" json:"provider_id"`
	Handle    string `db:"hash" json:"hash"`
}

func (db *DB) SaveBuild(b *Build) error {
	return db.SaveBuilds([]*Build{b})
}

func (db *DB) SaveBuilds(builds []*Build) error {
	fmt.Println("SAVING BUILDS: ", builds)
	return nil
}
