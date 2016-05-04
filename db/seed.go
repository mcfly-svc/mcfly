package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mikec/msplapi/models"
)

// Adds seed data to the database
func Seed(dbUrl string) {
	Clean(dbUrl)

	db, err := models.NewDB(dbUrl)
	checkFatal(err)

	for _, val := range []string{"jabroni.com", "schlockbox"} {
		addProviderValue(db.DB, val)
	}

	u := &models.User{Name: "Mock Johnson", AccessToken: "mock_access_token_123"}
	insertUser(db, u)

	insertProviderAccessToken(db, u.ID, &models.ProviderAccessToken{
		"jabroni.com",
		"mattmocks",
		"mock_jabroni.com_token_123",
	})

	insertProviderAccessToken(db, u.ID, &models.ProviderAccessToken{
		"schlockbox",
		"mattmocks@gmail.com",
		"mock_schlockbox_token_123",
	})

}

func addProviderValue(db *sqlx.DB, val string) {
	fmt.Printf("Add provider value \"%s\"\n", val)
	_, err := db.Exec(fmt.Sprintf("ALTER TYPE provider ADD VALUE '%s'", val))
	if !isDbErr(err, "duplicate_object") {
		checkFatal(err)
	}
}

func insertUser(db *models.DB, u *models.User) {
	fmt.Printf("INSERT marsupi_user %+v\n", *u)
	err := db.SaveUser(u)
	checkFatal(err)
}

func insertProviderAccessToken(db *models.DB, uid int64, pt *models.ProviderAccessToken) {
	fmt.Printf("INSERT provider_access_token %+v for user %d\n", *pt, uid)
	err := db.SetUserProviderToken(uid, pt)
	checkFatal(err)
}

/*func insertUser(
	db *sqlx.DB,
	name string,
	accessToken string,
) int64 {
	fmt.Printf("INSERT user name=%s, access_token=%s\n", name, accessToken)
	r := db.QueryRow(`INSERT INTO marsupi_user (name, access_token)
							 VALUES ($1, $2)
							 RETURNING id`, name, accessToken)
	var id int64
	if err := r.Scan(&id); err != nil {
		log.Fatalln(err)
	}
	return id
}

func insertProviderAccessToken(
	db *sqlx.DB,
	provider string,
	username string,
	accessToken string,
	userId int64,
) {
	fmt.Printf("INSERT provider_access_token provider=%s, access_token=%s\n", name, accessToken)
	r := db.Exec(`INSERT INTO provider_access_token (provider, provider_username, access_token, user_id)
								VALUES ($1, $2, $3, $4);`, name, accessToken)
	var id int64
	if err := r.Scan(&id); err != nil {
		log.Fatalln(err)
	}
	return id
}*/

/*
#!/bin/sh

echo "cleaning the Database"
msplapi db.clean

echo "change supported providers to test values"
psql -d marsupi_test <<- EOF
	ALTER TYPE provider ADD VALUE 'jabroni.com';
	ALTER TYPE provider ADD VALUE 'schlockbox';
EOF

echo "adding seed users"

psql -d marsupi_test <<- EOF

	INSERT INTO marsupi_user (name, access_token)
	VALUES ('Matt Mockerson', 'mock_token_123')
	RETURNING id;
	\gset
	\echo new user id=:id

	INSERT INTO provider_access_token (provider, provider_username, access_token, user_id)
	VALUES ('jabroni.com', 'mattmocks', 'mock_jabroni.com_token_123', :id);
	\echo new provider_access_token 'mock_jabroni.com_token_123'

	INSERT INTO provider_access_token (provider, provider_username, access_token, user_id)
	VALUES ('schlockbox', 'mattmocks@gmail.com', 'mock_schlockbox_token_123', :id);
	\echo new provider_access_token 'mock_schlockbox_token_123'

EOF

*/
