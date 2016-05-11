package models_test

import "testing"

func TestDeleteUserProject(t *testing.T) {
	resetDB()
	user, err := DB.GetUserByAccessToken("mock_seeded_access_token_123")
	checkErr(t, err)
	err = DB.DeleteUserProject(user, "jabroni.com", "mattmocks/project-2")
	checkErr(t, err)
	projects, err := DB.GetUserProjects(user)
	checkErr(t, err)
	if len(projects) != 2 {
		t.Error("Number of projects should be 2")
	}
	allProjects, err := DB.GetAllProjects()
	checkErr(t, err)
	if len(allProjects) != 2 {
		t.Error("Number of projects should be 2")
	}
}
