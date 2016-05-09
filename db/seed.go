package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mikec/msplapi/models"
)

// Seed inserts seed data into the database. It runs Clean first to delete any existing data.
func Seed(dbUrl string) {
	Clean(dbUrl)

	db, err := models.NewDB(dbUrl)
	checkFatal(err)

	for _, val := range []string{"jabroni.com", "schlockbox"} {
		addProviderValue(db.DB, val)
	}

	u := &models.User{Name: "Matt Mockman", AccessToken: "mock_seeded_access_token_123"}
	insertUser(db, u)

	insertProviderAccessToken(db, u.ID, &models.ProviderAccessToken{
		Provider:         "jabroni.com",
		ProviderUsername: "mattmocks",
		AccessToken:      "mock_jabroni.com_token_123",
	})

	insertProviderAccessToken(db, u.ID, &models.ProviderAccessToken{
		Provider:         "schlockbox",
		ProviderUsername: "mattmocks@gmail.com",
		AccessToken:      "mock_schlockbox_token_123",
	})

	u2 := &models.User{Name: "Penelope Providerless", AccessToken: "mock_token_for_user_with_no_provider_tokens"}
	insertUser(db, u2)

	u3 := &models.User{
		Name:        "Bethany Badprovidertoken",
		AccessToken: "mock_token_for_user_with_bad_jabroni.com_token",
	}
	insertUser(db, u3)

	insertProviderAccessToken(db, u3.ID, &models.ProviderAccessToken{
		Provider:         "jabroni.com",
		ProviderUsername: "bbadprovidertoken",
		AccessToken:      "bad_saved_jabroni.com_token_123",
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
