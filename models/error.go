package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type QueryExecError struct {
	Query   string
	DbError error
	Name    string
	Args    []interface{}
	NoRows  bool
	TxDone  bool
}

func (e QueryExecError) Error() string {
	if e.Args == nil {
		return fmt.Sprintf("Query '%s' failed with error: %v", e.Query, e.DbError)
	} else {
		return fmt.Sprintf("Query '%s' with args %q failed with error: %v", e.Query, e.Args, e.DbError)
	}
}

func NewQueryError(query string, queryError error, args []interface{}) *QueryExecError {
	err, ok := queryError.(*pq.Error)
	if !ok {
		noRows := queryError == sql.ErrNoRows
		txCommitted := queryError == sql.ErrTxDone
		return &QueryExecError{query, queryError, "", args, noRows, txCommitted}
	}
	return &QueryExecError{query, err, err.Code.Name(), args, false, false}
}
