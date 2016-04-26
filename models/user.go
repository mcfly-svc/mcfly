package models

import (
	"github.com/lib/pq"
)

type User struct {
	ID							int64			`db:"id" json:"id"`
	Name						string		`db:"name" json:"name"`
	GitHubToken			string		`db:"github_token" json:"github_token"`
	GitHubUsername	string		`db:"github_username" json:"github_username"`
}

func (db *DB) SaveUser(u *User) error {
	q := `INSERT INTO marsupi_user(name,github_token,github_username) VALUES($1,$2,$3) RETURNING id`
	r := db.QueryRow(q, u.Name, u.GitHubToken, u.GitHubUsername)

	var id int64
	if err := r.Scan(&id); err != nil {
		err, ok := err.(*pq.Error)
		if !ok {
			return err
		}
		return &QueryExecError{"SaveUser", q, err, err.Code.Name()}
	}

	u.ID = id

	return nil
}

func (db *DB) GetUserByGitHubToken(token string) (*User, error) {
	user := &User{}
	q := `SELECT * FROM marsupi_user WHERE github_token=$1`
	err := db.Get(user, q, token)
	if err != nil {
		err, ok := err.(*pq.Error)
		if !ok {
			return nil, err
		}
		return nil, &QueryExecError{"GetUserByGitHubToken", q, err, err.Code.Name()}
	}

	return user, nil
}

func (db *DB) GetUserById(id int64) (*User, error) {
	user := &User{}
	err := db.Get(user, `SELECT * FROM marsupi_user WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *DB) DeleteUser(id int64) error {
	q := `DELETE FROM marsupi_user WHERE id=$1`
	_, err := db.Exec(q, id)
	if err != nil {
		return &QueryExecError{"DeleteUser", q, err, ""}
	}
	return nil
}

func (db *DB) GetUsers() ([]User, error) {
	users := []User{}
	err := db.Select(&users, `SELECT * FROM marsupi_user`)
	if err != nil {
		return nil, err
	}
	return users, nil
}
