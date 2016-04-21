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

type EndpointTestRunner struct {
	Endpoint 			client.EntityEndpoint
	Entity1				string
	Entity2 			string
}

func (self *EndpointTestRunner) RunAll(t *testing.T) {
  self.RunCreateTest(t)
  self.RunGetAllTest(t)
  self.RunGetTest(t)
  self.RunMissingTest(t)
  self.RunDuplicateTest(t)
  self.RunCreateWithInvalidJsonTest(t)
  self.RunDeleteTest(t)
}

// create an entity
func (self *EndpointTestRunner) RunCreateTest(t *testing.T) {
	cleanupDB()

  ec := EndpointTestClient{t, self.Endpoint}
  e, res := ec.Create(self.Entity1)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)

  assert.Equal(t, e.ID > 0, true, fmt.Sprintf("Create %s did not return an ID", self.Endpoint.Name))
}

// get all entities
func (self *EndpointTestRunner) RunGetAllTest(t *testing.T) {
	cleanupDB()

  ec := EndpointTestClient{t, self.Endpoint}
  ec.Create(self.Entity1)
  ec.Create(self.Entity2)
  entities := ec.GetAllEntities()

  assert.Len(t, entities, 2, fmt.Sprintf("Wrong number of %s", self.Endpoint.Name))
}

// create an entity and get it by ID
func (self *EndpointTestRunner) RunGetTest(t *testing.T) {
  cleanupDB()

  ec := EndpointTestClient{t, self.Endpoint}
  ec.Create(self.Entity1)
  entities := ec.GetAllEntities()
  e := ec.GetEntity(entities[0].ID)

  assert.Equal(t, e.ID, entities[0].ID)
}

// get an entity that doesn't exist
func (self *EndpointTestRunner) RunMissingTest(t *testing.T) {
  cleanupDB()

  ec := EndpointTestClient{t, self.Endpoint}
  res := ec.Get(123)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.ApiError{
  	fmt.Sprintf("Failed to get %s where id=123", self.Endpoint.Name),
  })

}

// creating two duplicate entites should fail
func (self *EndpointTestRunner) RunDuplicateTest(t *testing.T) {
  cleanupDB()

  ec := EndpointTestClient{t, self.Endpoint}
  ec.Create(self.Entity1)
  _, res := ec.Create(self.Entity1)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.ApiError{
    fmt.Sprintf("%s already exists", util.Capitalize(self.Endpoint.SingularName)),
  })
}

// creating an entity with invalid json should fail
func (self *EndpointTestRunner) RunCreateWithInvalidJsonTest(t *testing.T) {
	cleanupDB()

  ec := EndpointTestClient{t, self.Endpoint}
	_, res := ec.Create(`{ "bad" }`)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.InvalidJsonApiErr)
}

// creating an entity, then deleting it, should return 200 status and delete the entity
func (self *EndpointTestRunner) RunDeleteTest(t *testing.T) {
  cleanupDB()

  ec := EndpointTestClient{t, self.Endpoint}
  ec.Create(self.Entity1)
  entities := ec.GetAllEntities()
  res := ec.Delete(entities[0].ID)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)

  entities = ec.GetAllEntities()

  assert.Len(t, entities, 0)
}
