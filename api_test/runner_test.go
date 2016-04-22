package api_test

import (
	"github.com/mikec/marsupi-api/api"
	"github.com/mikec/marsupi-api/client"
  "github.com/mikec/marsupi-api/util"

  "github.com/stretchr/testify/assert"

  "fmt"
	"testing"
  "reflect"
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
  res := c.Create(self.Entity1)
  d := dataResp{res}
  id := d.ID()

  assertStatusCode(t, 200, res.StatusCode, self.Endpoint.Name)
  assert.Equal(t, true, id > 0, fmt.Sprintf("Create %s did not return an ID", self.Endpoint.Name))
}


// get all entities
func (self *Runner) RunGetAllTest(t *testing.T) {
	cleanupDB()

  c := Client{t, self.Endpoint}
  c.Create(self.Entity1)
  c.Create(self.Entity2)
  res := c.GetAll()
  d := dataArrayResp{res}
  n := d.Len()

  assertStatusCode(t, 200, res.StatusCode, self.Endpoint.Name)
  assert.Equal(t, 2, n, fmt.Sprintf("GetAll %s returned wrong number of %s", self.Endpoint.Name, self.Endpoint.Name))
}

// create an entity and get it by ID
func (self *Runner) RunGetTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  res := c.Create(self.Entity1)
  d := dataResp{res}
  id1 := d.ID()

  res = c.Get(id1)
  d = dataResp{res}
  id2 := d.ID()

  assertStatusCode(t, 200, res.StatusCode, self.Endpoint.Name)
  assert.Equal(t, id1, id2, fmt.Sprintf("Get %s failed, id value mismatched", self.Endpoint.SingularName))
}


// get an entity that doesn't exist
func (self *Runner) RunMissingTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  res := c.Get(123)

  assertStatusCode(t, 400, res.StatusCode, self.Endpoint.Name)

  apiErr := api.ApiError{fmt.Sprintf("Failed to get %s where id=123", self.Endpoint.Name)}
  assert.Equal(t, apiErr, res.Data)
}

// creating two duplicate entites should fail
func (self *Runner) RunDuplicateTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  c.Create(self.Entity1)
  res := c.Create(self.Entity1)

  assertStatusCode(t, 400, res.StatusCode, self.Endpoint.Name)

  apiErr := api.ApiError{fmt.Sprintf("%s already exists", util.Capitalize(self.Endpoint.SingularName))}
  assert.Equal(t, apiErr, res.Data)
}

// creating an entity with invalid json should fail
func (self *Runner) RunCreateWithInvalidJsonTest(t *testing.T) {
	cleanupDB()

  c := Client{t, self.Endpoint}
	res := c.Create(`{ "bad" }`)

  assertStatusCode(t, 400, res.StatusCode, self.Endpoint.Name)
  assert.Equal(t, api.InvalidJsonApiErr, res.Data)
}

// creating an entity, then deleting it, should return 200 status and delete the entity
func (self *Runner) RunDeleteTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  c.Create(self.Entity1)
  res := c.GetAll()
  d := dataArrayResp{res}

  res = c.Delete(d.FirstItemID())

  assertStatusCode(t, 200, res.StatusCode, self.Endpoint.Name)

  res = c.GetAll()
  d = dataArrayResp{res}
  n := d.Len()

  assert.Equal(t, 0, n, fmt.Sprintf("Expected 0 %s after Delete, but got %d", self.Endpoint.Name, n))
}


type dataResp struct {
  *client.ClientResponse
}

func (d dataResp) ID() int64 {
  return id(elem(d.Data))
}

type dataArrayResp struct {
  *client.ClientResponse
}

func (d dataArrayResp) Len() int {
  return elem(d.Data).Len()
}

func (d dataArrayResp) FirstItemID() int64 {
  return id(elem(d.Data).Index(0))
}

func id(v reflect.Value) int64 {
  return v.FieldByName("ID").Interface().(int64)
}

func elem(v interface{}) reflect.Value {
  return reflect.ValueOf(v).Elem()
}

func assertStatusCode(t *testing.T, actualCode int, expectCode int, endpointName string) {
  assert.Equal(t, actualCode, expectCode, fmt.Sprintf("Create %s returned the wrong status code", endpointName))
}
