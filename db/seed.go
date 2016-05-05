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

	u := &models.User{Name: "Matt Mockman", AccessToken: "mock_seeded_access_token_123"}
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
