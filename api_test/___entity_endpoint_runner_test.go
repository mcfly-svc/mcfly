package api_test

/*
import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mikec/marsupi-api/api"
	"github.com/mikec/marsupi-api/client"
	"github.com/stretchr/testify/assert"
)

type EntityEndpointRunner struct {
	Endpoint client.EntityEndpoint
	Entity1  string
	Entity2  string
}

// create an entity
func (self *EntityEndpointRunner) RunCreateTest(t *testing.T) {
	cleanupDB()

	c := EndpointClient{t, self.Endpoint}
	res, _ := c.Create(self.Entity1)
	d := dataResp{res}
	id := d.ID()

	assertStatusCode(t, 200, res.StatusCode)
	assert.Equal(t, true, id > 0, fmt.Sprintf("Create %s did not return an ID"))
}

// get all entities
func (self *EntityEndpointRunner) RunGetAllTest(t *testing.T) {
	cleanupDB()

	c := EndpointClient{t, self.Endpoint}
	c.Create(self.Entity1)
	c.Create(self.Entity2)
	res, _ := c.GetAll()
	d := dataArrayResp{res}
	n := d.Len()

	assertStatusCode(t, 200, res.StatusCode)
	assert.Equal(t, 2, n, fmt.Sprintf("GetAll %s returned wrong number of %s", self.Endpoint.Name, self.Endpoint.Name))
}

// create an entity and get it by ID
func (self *EntityEndpointRunner) RunGetTest(t *testing.T) {
	cleanupDB()

	c := EndpointClient{t, self.Endpoint}
	res, _ := c.Create(self.Entity1)

	d := dataResp{res}
	id1 := d.ID()

	res, httpRes := c.Get(id1)
	d = dataResp{res}
	id2 := d.ID()

	assertStatusCode(t, 200, res.StatusCode)
	assert.Equal(t, id1, id2, fmt.Sprintf("Get %s failed, id value mismatched", self.Endpoint.SingularName))

	assertJsonEncoded(t, httpRes, "id")
}

// get an entity that doesn't exist
func (self *EntityEndpointRunner) RunMissingTest(t *testing.T) {
	cleanupDB()

	c := EndpointClient{t, self.Endpoint}
	res, _ := c.Get(123)

	assertStatusCode(t, 400, res.StatusCode)
	assert.Equal(t, *api.NewGetEntityErr(self.Endpoint.SingularName, 123), res.Data)
}

// creating two duplicate entites should fail
func (self *EntityEndpointRunner) RunDuplicateTest(t *testing.T) {
	cleanupDB()

	c := EndpointClient{t, self.Endpoint}
	c.Create(self.Entity1)
	res, _ := c.Create(self.Entity1)

	assertStatusCode(t, 400, res.StatusCode)
	assert.Equal(t, *api.NewDuplicateCreateErr(self.Endpoint.SingularName), res.Data)
}

// creating an entity with invalid json should fail
func (self *EntityEndpointRunner) RunCreateWithInvalidJsonTest(t *testing.T) {
	cleanupDB()

	c := EndpointClient{t, self.Endpoint}
	res, _ := c.Create(`{ "bad" }`)

	assertStatusCode(t, 400, res.StatusCode)
	assert.Equal(t, *api.NewInvalidJsonErr(), res.Data)
}

// creating an entity, then deleting it, should return 200 status and delete the entity
func (self *EntityEndpointRunner) RunDeleteTest(t *testing.T) {
	cleanupDB()

	c := EndpointClient{t, self.Endpoint}
	res, _ := c.Create(self.Entity1)
	d := dataResp{res}
	id := d.ID()

	res, _ = c.Delete(id)

	assertStatusCode(t, 200, res.StatusCode)

	res, _ = c.Get(id)

	assertStatusCode(t, 400, res.StatusCode)
	assert.Equal(t, *api.NewGetEntityErr(self.Endpoint.SingularName, id), res.Data)
}

// try to get an entity using an invalid ID
func (self *EntityEndpointRunner) RunInvalidGetTest(t *testing.T) {
	cleanupDB()

	c := EndpointClient{t, self.Endpoint}
	res, _ := c.Get(0)

	assertStatusCode(t, 400, res.StatusCode)
	assert.Equal(t, *api.NewInvalidParamErr("ID", "0"), res.Data)
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
*/
