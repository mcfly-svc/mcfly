package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestGetProject(t *testing.T) {
	resetDB()
	p, err := DB.GetProject("mattmocks/project-2", "jabroni.com")
	checkErr(t, err)
	assert.NotZero(t, p.ID, "GetProject should return a project with a non-zero ID")
	assert.Equal(t, "mattmocks/project-2", p.Handle, "GetProject should return a project with the given handle")
}

func TestGetProjectThatDoesNotExist(t *testing.T) {
	resetDB()
	p, err := DB.GetProject("project/that-does-not-exist", "jabroni.com")
	checkErr(t, err)
	assert.Nil(t, p, "GetProject should return nil")
}
