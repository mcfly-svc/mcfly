package api_test

import (
	"github.com/mikec/marsupi-api/api"
	"github.com/mikec/marsupi-api/testutil"
	"github.com/mikec/marsupi-api/apiutil"

  "github.com/stretchr/testify/assert"

  "fmt"
  "reflect"
	"testing"
)

type EndpointTests struct {
	Endpoint 			apiutil.EntityEndpoint
	Entity1				string
	Entity2 			string
}

// create an entity
func (self *EndpointTests) RunCreateTest(t *testing.T) {
	cleanupDB()

  et := EndpointTester{t, self.Endpoint}
  res := et.Create(self.Entity1)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)
}

// get all entities
func (self *EndpointTests) RunGetAllTest(t *testing.T) {
	cleanupDB()

  et := EndpointTester{t, self.Endpoint}
  et.Create(self.Entity1)
  et.Create(self.Entity2)
  entities := et.GetAll()

  assert.Len(t, entities, 2, fmt.Sprintf("Wrong number of %s", self.Endpoint.Name))
}

// create an entity and get it by ID
func (self *EndpointTests) RunGetTest(t *testing.T) {
  cleanupDB()

  et := EndpointTester{t, self.Endpoint}
  et.Create(self.Entity1)
  entities := et.GetAll()
  id1 := reflect.ValueOf(entities[0]).FieldByName("ID").Interface().(int64)
  e := et.Get(id1)
  id2 := reflect.ValueOf(e).FieldByName("ID").Interface().(int64)

  assert.Equal(t, id2, id1)
}

// get an entity that doesn't exist
func (self *EndpointTests) RunMissingTest(t *testing.T) {
  cleanupDB()

  et := EndpointTester{t, self.Endpoint}
  res := et.GetRes(123)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.ApiError{
  	fmt.Sprintf("Failed to get %s where id=123", self.Endpoint.Name),
  })

}

// creating two duplicate entites should fail
func (self *EndpointTests) RunDuplicateTest(t *testing.T) {
  cleanupDB()

  et := EndpointTester{t, self.Endpoint}
  et.Create(self.Entity1)
  res := et.Create(self.Entity1)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
}

// creating an entity with invalid json should fail
func (self *EndpointTests) RunCreateWithInvalidJsonTest(t *testing.T) {
	cleanupDB()

  et := EndpointTester{t, self.Endpoint}
	res := et.Create(`{ "bad" }`)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.InvalidJsonApiErr)
}

// creating an entity, then deleting it, should return 200 status and delete the entity
func (self *EndpointTests) RunDeleteTest(t *testing.T) {
  cleanupDB()

  et := EndpointTester{t, self.Endpoint}
  et.Create(self.Entity1)
  entities := et.GetAll()
  id1 := reflect.ValueOf(entities[0]).FieldByName("ID").Interface().(int64)
  res := et.Delete(id1)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)

  entities = et.GetAll()

  assert.Len(t, entities, 0)
}
