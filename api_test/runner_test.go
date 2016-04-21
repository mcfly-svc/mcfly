package api_test

import (
	"github.com/mikec/marsupi-api/api"
	"github.com/mikec/marsupi-api/testutil"
	"github.com/mikec/marsupi-api/client"
  "github.com/mikec/marsupi-api/util"

  "github.com/stretchr/testify/assert"

  "fmt"
	"testing"
)

type Runner struct {
	Endpoint 			client.EntityEndpoint
	Entity1				string
	Entity2 			string
}

// create an entity
func (self *Runner) RunCreateTest(t *testing.T) {
	cleanupDB()

  c := Client{t, self.Endpoint}
  e, res := c.Create(self.Entity1)

  rt := testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)

  assert.Equal(t, e.ID > 0, true, fmt.Sprintf("Create %s did not return an ID", self.Endpoint.Name))
}

// get all entities
func (self *Runner) RunGetAllTest(t *testing.T) {
	cleanupDB()

  c := Client{t, self.Endpoint}
  c.Create(self.Entity1)
  c.Create(self.Entity2)
  entities, res := c.GetAll()

  rt := testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)

  assert.Len(t, entities, 2, fmt.Sprintf("Wrong number of %s", self.Endpoint.Name))
}

// create an entity and get it by ID
func (self *Runner) RunGetTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  c.Create(self.Entity1)
  entities, _ := c.GetAll()
  e, res := c.Get(entities[0].ID)

  rt := testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)

  assert.Equal(t, e.ID, entities[0].ID)
}

// get an entity that doesn't exist
func (self *Runner) RunMissingTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  _, res := c.Get(123)

  rt := testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.ApiError{
  	fmt.Sprintf("Failed to get %s where id=123", self.Endpoint.Name),
  })
}

// creating two duplicate entites should fail
func (self *Runner) RunDuplicateTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  c.Create(self.Entity1)
  _, res := c.Create(self.Entity1)

  rt := testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.ApiError{
    fmt.Sprintf("%s already exists", util.Capitalize(self.Endpoint.SingularName)),
  })
}

// creating an entity with invalid json should fail
func (self *Runner) RunCreateWithInvalidJsonTest(t *testing.T) {
	cleanupDB()

  c := Client{t, self.Endpoint}
	_, res := c.Create(`{ "bad" }`)

  rt := testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.InvalidJsonApiErr)
}

// creating an entity, then deleting it, should return 200 status and delete the entity
func (self *Runner) RunDeleteTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  c.Create(self.Entity1)
  entities, _ := c.GetAll()
  res := c.Delete(entities[0].ID)

  rt := testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)

  entities, _ = c.GetAll()

  assert.Len(t, entities, 0)
}
