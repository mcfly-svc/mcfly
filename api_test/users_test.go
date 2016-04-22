package api_test

import (
	"testing"
)

var usersRunner Runner

func init() {
  usersRunner = Runner{
    apiClient.Users,
    `{"github_token": "a1b2c3"}`,
    `{"github_token": "a1b2c3"}`,
  }
}

func TestCreateUser(t *testing.T) { usersRunner.RunCreateTest(t) }
func TestGetUser(t *testing.T) { usersRunner.RunGetTest(t) }
func TestGetMissingUser(t *testing.T) { usersRunner.RunMissingTest(t) }
func TestCreateDuplicateUser(t *testing.T) { usersRunner.RunDuplicateTest(t) }
func TestCreateUserWithInvalidJson(t *testing.T) { usersRunner.RunCreateWithInvalidJsonTest(t) }
func TestDeleteUser(t *testing.T) { usersRunner.RunDeleteTest(t) }
func TestGetUserInvalidID(t *testing.T) { usersRunner.RunInvalidGetTest(t) }

// no get all endpoint for users yet
//func TestGetAllUsers(t *testing.T) { usersRunner.RunGetAllTest(t) }