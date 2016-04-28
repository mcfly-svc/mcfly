package models

type User struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	AccessToken string `db:"access_token" json:"access_token"`
}

type ProviderAccessToken struct {
	Provider         string `db:"provider" json:"provider"`
	ProviderUsername string `db:"provider_username" json:"provider_username"`
	AccessToken      string `db:"access_token" json:"access_token"`
}

func (db *DB) SaveUser(u *User) error {
	q := `INSERT INTO marsupi_user(name, access_token) VALUES($1, $2) RETURNING id`
	r := db.QueryRow(q, u.Name, u.AccessToken)

	var id int64
	if err := r.Scan(&id); err != nil {
		return handleQueryError("SaveUser", q, err)
	}

	u.ID = id

	return nil
}

func (db *DB) SetUserProviderToken(userID int64, providerToken ProviderAccessToken) error {
	q := `INSERT INTO provider_access_token(user_id, provider, provider_username, access_token)
				VALUES($1,$2,$3,$4)`
	_, err := db.Exec(q, userID, providerToken.Provider, providerToken.ProviderUsername, providerToken.AccessToken)
	if err != nil {
		return handleQueryError("SetUserProviderToken", q, err)
	}
	return nil
}

func (db *DB) GetUserByProviderToken(providerToken ProviderAccessToken) (*User, error) {
	user := &User{}
	q := `SELECT * FROM marsupi_user 
				INNER JOIN provider_access_token
				ON marsupi_user.id = provider_access_token.user_id
				WHERE provider_access_token.provider=$1
				AND provider_access_token.provider_username=$2
				AND provider_access_token.access_token=$3`
	err := db.Get(user, q, providerToken.Provider, providerToken.ProviderUsername, providerToken.AccessToken)
	if err != nil {
		return nil, handleQueryError("GetUserByProviderToken", q, err)
	}

	return user, nil
}

/*

GARBAGE

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
		return handleQueryError("DeleteUser", q, err)
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

*/
