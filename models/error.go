package models

import (
	"fmt"
)

type QueryExecError struct {
	OperationName 	string
	Query 					string
	DbError 				error
	Name 						string
}

func (e QueryExecError) Error() string {
	return fmt.Sprintf("%s query '%s' failed with error: %v", e.OperationName, e.Query, e.DbError)
}
