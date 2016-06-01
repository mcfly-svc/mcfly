package db

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mcfly-svc/mcfly/models"
)

// Seed inserts seed data into the database. It runs Clean first to delete any existing data.
func (mdb *McflyDB) Seed() {
	mdb.Clean()

	for _, val := range []string{"jabroni.com", "schlockbox"} {
		addProviderValue(mdb, val)
	}

	u := &models.User{Name: strPtr("Matt Mockman"), AccessToken: "mock_seeded_access_token_123"}
	insertUser(mdb, u)

	insertProviderAccessToken(mdb, u.ID, &models.ProviderAccessToken{
		Provider:         "jabroni.com",
		ProviderUsername: "mattmocks",
		AccessToken:      "mock_jabroni.com_token_123",
	})

	insertProviderAccessToken(mdb, u.ID, &models.ProviderAccessToken{
		Provider:         "schlockbox",
		ProviderUsername: "mattmocks@gmail.com",
		AccessToken:      "mock_schlockbox_token_123",
	})

	p1 := &models.Project{
		Handle:         "mattmocks/project-1",
		SourceProvider: "jabroni.com",
		SourceUrl:      "https://jabroni.com/mattmocks/project-1",
	}
	insertProject(mdb, u, p1)

	p2 := &models.Project{
		Handle:         "mattmocks/project-2",
		SourceProvider: "jabroni.com",
		SourceUrl:      "https://jabroni.com/mattmocks/project-2",
	}
	insertProject(mdb, u, p2)

	p3 := &models.Project{
		Handle:         "mattmocks/project-3",
		SourceProvider: "jabroni.com",
		SourceUrl:      "https://jabroni.com/mattmocks/project-3",
	}
	insertProject(mdb, u, p3)

	insertBuild(mdb, p3, &models.Build{
		Handle:       "abc-1",
		DeployStatus: "succeeded",
		ProviderUrl:  strPtr("https://jabroni.com/mattmocks/project-3/builds/abc-1"),
	})

	insertBuild(mdb, p3, &models.Build{
		Handle:       "abc-2",
		DeployStatus: "succeeded",
		ProviderUrl:  strPtr("https://jabroni.com/mattmocks/project-3/builds/abc-2"),
	})

	insertBuild(mdb, p3, &models.Build{
		Handle:       "abc-3",
		DeployStatus: "failed",
		ProviderUrl:  strPtr("https://jabroni.com/mattmocks/project-3/builds/abc-3"),
	})

	insertBuild(mdb, p3, &models.Build{
		Handle:       "abc-4",
		DeployStatus: "succeeded",
	})

	insertBuild(mdb, p3, &models.Build{
		Handle:       "abc-4",
		DeployStatus: "pending",
		ProviderUrl:  strPtr("https://jabroni.com/mattmocks/project-3/builds/abc-4"),
	})

	u2 := &models.User{Name: strPtr("Penelope Providerless"), AccessToken: "mock_token_for_user_with_no_provider_tokens"}
	insertUser(mdb, u2)

	u3 := &models.User{
		Name:        strPtr("Bethany Badprovidertoken"),
		AccessToken: "mock_token_for_user_with_bad_jabroni.com_token",
	}
	insertUser(mdb, u3)

	insertProviderAccessToken(mdb, u3.ID, &models.ProviderAccessToken{
		Provider:         "jabroni.com",
		ProviderUsername: "bbadprovidertoken",
		AccessToken:      "bad_saved_jabroni.com_token_123",
	})

}

func addProviderValue(mdb *McflyDB, val string) {
	//fmt.Printf("Add provider value \"%s\"\n", val)
	_, err := mdb.Exec(fmt.Sprintf("ALTER TYPE provider ADD VALUE '%s'", val))
	if !isDbErr(err, "duplicate_object") {
		checkDbNotFoundErr(err, mdb.DatabaseName)
	}
}

func insertUser(mdb *McflyDB, u *models.User) {
	//fmt.Printf("INSERT mcfly_user %+v\n", *u)
	err := mdb.SaveUser(u)
	checkDbNotFoundErr(err, mdb.DatabaseName)
}

func insertProviderAccessToken(mdb *McflyDB, uid int64, pt *models.ProviderAccessToken) {
	//fmt.Printf("INSERT provider_access_token %+v for user %d\n", *pt, uid)
	err := mdb.SetUserProviderToken(uid, pt)
	checkDbNotFoundErr(err, mdb.DatabaseName)
}

func insertProject(mdb *McflyDB, u *models.User, p *models.Project) {
	//fmt.Printf("INSERT project %+v\n for user %d\n", *p, u.ID)
	err := mdb.SaveProject(p, u)
	checkDbNotFoundErr(err, mdb.DatabaseName)
}

func insertBuild(mdb *McflyDB, p *models.Project, b *models.Build) {
	//fmt.Printf("INSERT build %+v\n for project %d\n", *b, p.ID)
	err := mdb.SaveBuild(b, p)
	checkDbNotFoundErr(err, mdb.DatabaseName)
}

func strPtr(s string) *string {
	return &s
}
