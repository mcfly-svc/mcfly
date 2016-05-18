package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mikec/msplapi/models"
)

// Seed inserts seed data into the database. It runs Clean first to delete any existing data.
func Seed(conn *sqlx.DB) {
	Clean(conn)
	db := &models.DB{conn}

	for _, val := range []string{"jabroni.com", "schlockbox"} {
		addProviderValue(db.DB, val)
	}

	u := &models.User{Name: strPtr("Matt Mockman"), AccessToken: "mock_seeded_access_token_123"}
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

	p1 := &models.Project{
		Handle:         "mattmocks/project-1",
		SourceProvider: "jabroni.com",
		SourceUrl:      "https://jabroni.com/mattmocks/project-1",
	}
	insertProject(db, u, p1)

	p2 := &models.Project{
		Handle:         "mattmocks/project-2",
		SourceProvider: "jabroni.com",
		SourceUrl:      "https://jabroni.com/mattmocks/project-2",
	}
	insertProject(db, u, p2)

	p3 := &models.Project{
		Handle:         "mattmocks/project-3",
		SourceProvider: "jabroni.com",
		SourceUrl:      "https://jabroni.com/mattmocks/project-3",
	}
	insertProject(db, u, p3)

	insertBuild(db, p3, &models.Build{
		Handle:       "abc-1",
		DeployStatus: "succeeded",
		ProviderUrl:  strPtr("https://jabroni.com/mattmocks/project-3/builds/abc-1"),
	})

	insertBuild(db, p3, &models.Build{
		Handle:       "abc-2",
		DeployStatus: "succeeded",
		ProviderUrl:  strPtr("https://jabroni.com/mattmocks/project-3/builds/abc-2"),
	})

	insertBuild(db, p3, &models.Build{
		Handle:       "abc-3",
		DeployStatus: "failed",
		ProviderUrl:  strPtr("https://jabroni.com/mattmocks/project-3/builds/abc-3"),
	})

	insertBuild(db, p3, &models.Build{
		Handle:       "abc-4",
		DeployStatus: "succeeded",
	})

	insertBuild(db, p3, &models.Build{
		Handle:       "abc-4",
		DeployStatus: "pending",
		ProviderUrl:  strPtr("https://jabroni.com/mattmocks/project-3/builds/abc-4"),
	})

	u2 := &models.User{Name: strPtr("Penelope Providerless"), AccessToken: "mock_token_for_user_with_no_provider_tokens"}
	insertUser(db, u2)

	u3 := &models.User{
		Name:        strPtr("Bethany Badprovidertoken"),
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
	//fmt.Printf("Add provider value \"%s\"\n", val)
	_, err := db.Exec(fmt.Sprintf("ALTER TYPE provider ADD VALUE '%s'", val))
	if !isDbErr(err, "duplicate_object") {
		checkFatal(err)
	}
}

func insertUser(db *models.DB, u *models.User) {
	//fmt.Printf("INSERT marsupi_user %+v\n", *u)
	err := db.SaveUser(u)
	checkFatal(err)
}

func insertProviderAccessToken(db *models.DB, uid int64, pt *models.ProviderAccessToken) {
	//fmt.Printf("INSERT provider_access_token %+v for user %d\n", *pt, uid)
	err := db.SetUserProviderToken(uid, pt)
	checkFatal(err)
}

func insertProject(db *models.DB, u *models.User, p *models.Project) {
	//fmt.Printf("INSERT project %+v\n for user %d\n", *p, u.ID)
	err := db.SaveProject(p, u)
	checkFatal(err)
}

func insertBuild(db *models.DB, p *models.Project, b *models.Build) {
	//fmt.Printf("INSERT build %+v\n for project %d\n", *b, p.ID)
	err := db.SaveBuild(b, p)
	checkFatal(err)
}

func strPtr(s string) *string {
	return &s
}
