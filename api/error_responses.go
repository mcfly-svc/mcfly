package api

import "fmt"

var (
	InvalidJsonApiErr = ApiError{"Invalid JSON"}
)

type ApiError struct {
  Error         string        `json:"error"`
}

func (self *ApiError) InvalidParam(name string, val string) {
	self.Error = fmt.Sprintf("Invalid parameter %s=%s", name, val)
}
