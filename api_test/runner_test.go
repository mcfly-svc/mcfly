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
  res, _ := c.Create(self.Entity1)
  d := dataResp{res}
  id := d.ID()

  assertStatusCode(t, 200, res.StatusCode)
  assert.Equal(t, true, id > 0, fmt.Sprintf("Create %s did not return an ID"))
}


// get all entities
func (self *Runner) RunGetAllTest(t *testing.T) {
	cleanupDB()

  c := Client{t, self.Endpoint}
  c.Create(self.Entity1)
  c.Create(self.Entity2)
  res, _ := c.GetAll()
  d := dataArrayResp{res}
  n := d.Len()

  assertStatusCode(t, 200, res.StatusCode)
  assert.Equal(t, 2, n, fmt.Sprintf("GetAll %s returned wrong number of %s", self.Endpoint.Name, self.Endpoint.Name))
}

// create an entity and get it by ID
func (self *Runner) RunGetTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  res, _ := c.Create(self.Entity1)

  d := dataResp{res}
  id1 := d.ID()

  res, _ = c.Get(id1)
  d = dataResp{res}
  id2 := d.ID()

  fmt.Printf("GOT:%+v\n", res.Data)

  assertStatusCode(t, 200, res.StatusCode)
  assert.Equal(t, id1, id2, fmt.Sprintf("Get %s failed, id value mismatched", self.Endpoint.SingularName))
}


// get an entity that doesn't exist
func (self *Runner) RunMissingTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  res, _ := c.Get(123)

  assertStatusCode(t, 400, res.StatusCode)

  apiErr := api.ApiError{fmt.Sprintf("Failed to get %s where id=123", self.Endpoint.SingularName)}
  assert.Equal(t, apiErr, res.Data)
}

// creating two duplicate entites should fail
func (self *Runner) RunDuplicateTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  c.Create(self.Entity1)
  res, _ := c.Create(self.Entity1)

  assertStatusCode(t, 400, res.StatusCode)

  apiErr := api.ApiError{fmt.Sprintf("%s already exists", util.Capitalize(self.Endpoint.SingularName))}
  assert.Equal(t, apiErr, res.Data)
}

// creating an entity with invalid json should fail
func (self *Runner) RunCreateWithInvalidJsonTest(t *testing.T) {
	cleanupDB()

  c := Client{t, self.Endpoint}
	res, _ := c.Create(`{ "bad" }`)

  assertStatusCode(t, 400, res.StatusCode)
  assert.Equal(t, api.InvalidJsonApiErr, res.Data)
}

// creating an entity, then deleting it, should return 200 status and delete the entity
func (self *Runner) RunDeleteTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  res, _ := c.Create(self.Entity1)
  d := dataResp{res}
  id := d.ID()

  res, _ = c.Delete(id)

  assertStatusCode(t, 200, res.StatusCode)

  res, _ = c.Get(id)

  assertStatusCode(t, 400, res.StatusCode)
  apiErr := api.ApiError{fmt.Sprintf("Failed to get %s where id=%d", self.Endpoint.SingularName, id)}
  assert.Equal(t, apiErr, res.Data)
}

// try to get an entity using an invalid ID
func (self *Runner) RunInvalidGetTest(t *testing.T) {
  cleanupDB()

  c := Client{t, self.Endpoint}
  res, _ := c.Get(0)

  assertStatusCode(t, 400, res.StatusCode)

  apiErr := &api.ApiError{}
  apiErr.InvalidParam("ID", "0")
  assert.Equal(t, *apiErr, res.Data)
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

func assertStatusCode(t *testing.T, actualCode int, expectCode int) {
  assert.Equal(t, actualCode, expectCode, "Wrong status code")
}
