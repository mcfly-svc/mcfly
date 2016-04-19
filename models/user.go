package models

import "fmt"

type User struct {
	ID							int64			`db:"id" 									json:"id"`
	Name						string		`db:"name" 								json:"name"`
	GitHubToken			string		`db:"github_token" 				json:"github_token"`
	GitHubUsername	string		`db:"github_username" 		json:"github_username"`
}

func (db *DB) SaveUser(u *User) (*User, error) {
	var uid int64
	q := `INSERT INTO marsupi_user(name,github_token,github_username) VALUES($1,$2,$3) RETURNING id`
	r := db.QueryRow(q, u.Name, u.GitHubToken, u.GitHubUsername)

	err := r.Scan(&uid)
	if err != nil {
		return nil, &QueryExecError{"SaveUser", q, err}
	}

	newUser, err := db.GetUser(uid)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (db *DB) GetUser(id int64) (*User, error) {
	user := &User{}
	err := db.Get(user, `SELECT * FROM marsupi_user`)
	if err != nil {
		fmt.Println("ERR", err)
		return nil, err
	}
	return user, nil
}
