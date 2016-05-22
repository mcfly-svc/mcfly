package models_test

import (
	"testing"

	"github.com/mikec/msplapi/models"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUserProject(t *testing.T) {
	resetDB()
	user, err := mdb.GetUserByAccessToken("mock_seeded_access_token_123")
	checkErr(t, err)
	err = mdb.DeleteUserProject(user, "jabroni.com", "mattmocks/project-2")
	checkErr(t, err)
	projects, err := mdb.GetUserProjects(user)
	checkErr(t, err)
	if len(projects) != 2 {
		t.Error("Number of projects should be 2")
	}
	allProjects, err := mdb.GetAllProjects()
	checkErr(t, err)
	if len(allProjects) != 2 {
		t.Error("Number of projects should be 2")
	}
}

func TestGetProject(t *testing.T) {
	resetDB()
	p, err := mdb.GetProject("mattmocks/project-2", "jabroni.com")
	checkErr(t, err)
	assert.NotZero(t, p.ID, "GetProject should return a project with a non-zero ID")
	assert.Equal(t, "mattmocks/project-2", p.Handle, "GetProject should return a project with the given handle")
}

func TestGetProjectThatDoesNotExist(t *testing.T) {
	resetDB()
	p, err := mdb.GetProject("project/that-does-not-exist", "jabroni.com")
	assert.Nil(t, p, "GetProject should return nil")
	assert.Equal(t, err, models.ErrNotFound)
}
