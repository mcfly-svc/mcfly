package models

type User struct {
	ID          int64   `db:"id" json:"id"`
	Name        *string `db:"name" json:"name"`
	AccessToken string  `db:"access_token" json:"access_token"`
}

type ProviderAccessToken struct {
	Provider         string `db:"provider" json:"provider"`
	ProviderUsername string `db:"provider_username" json:"provider_username"`
	AccessToken      string `db:"access_token" json:"access_token"`
	UserID           int64  `db:"user_id" json:"user_id"`
}

// SaveUser saves a user and returns the new user's ID
func (db *DB) SaveUser(u *User) error {
	var id int64
	q := `INSERT INTO marsupi_user(name, access_token) VALUES($1, $2) RETURNING id`
	err := db.QueryRowScan(&id, q, u.Name, u.AccessToken)
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

// SetUserProviderToken sets a provider access token for a user
func (db *DB) SetUserProviderToken(userID int64, providerToken *ProviderAccessToken) error {
	q := `INSERT INTO provider_access_token(user_id, provider, provider_username, access_token)
				VALUES($1,$2,$3,$4)`
	_, err := db.Exec(q, userID, providerToken.Provider, providerToken.ProviderUsername, providerToken.AccessToken)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByAccessToken gets a user by their Marsupi access token
func (db *DB) GetUserByAccessToken(accessToken string) (*User, error) {
	user := &User{}
	q := `SELECT * FROM marsupi_user WHERE access_token=$1`
	err := db.Get(user, q, accessToken)
	if err != nil {
		if err.NoRows {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return user, nil
}

// GetUserByProviderToken gets a user by one of their provider access tokens
func (db *DB) GetUserByProviderToken(providerToken *ProviderAccessToken) (*User, error) {
	user := &User{}
	q := `SELECT marsupi_user.* FROM marsupi_user 
				INNER JOIN provider_access_token
				ON marsupi_user.id = provider_access_token.user_id
				WHERE provider_access_token.provider=$1
				AND provider_access_token.provider_username=$2
				AND provider_access_token.access_token=$3`
	err := db.Get(user, q, providerToken.Provider, providerToken.ProviderUsername, providerToken.AccessToken)
	if err != nil {
		if err.NoRows {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return user, nil
}

// GetUserByProject gets a user by one of their project handles
func (db *DB) GetUserByProject(project *Project) (*User, error) {
	user := &User{}
	q := `SELECT marsupi_user.* FROM marsupi_user 
				INNER JOIN user_project ON marsupi_user.id = user_project.user_id
				INNER JOIN project ON user_project.project_id = project.id
				WHERE project.handle=$1
				AND project.source_provider=$2`
	err := db.Get(user, q, project.Handle, project.SourceProvider)
	if err != nil {
		if err.NoRows {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return user, nil
}

// GetProviderTokenForUser gets a provider token for a given user and provider
func (db *DB) GetProviderTokenForUser(user *User, provider string) (*ProviderAccessToken, error) {
	pt := &ProviderAccessToken{}
	q := `SELECT provider_access_token.* FROM marsupi_user
				INNER JOIN provider_access_token
				ON marsupi_user.id = provider_access_token.user_id
				WHERE provider_access_token.provider=$1
				AND marsupi_user.id=$2`
	err := db.Get(pt, q, provider, user.ID)
	if err != nil {
		if err.NoRows {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return pt, nil
}
