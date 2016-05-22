package models_test

import (
	"testing"

	"github.com/mikec/msplapi/models"
	"github.com/stretchr/testify/assert"
)

func TestConnectionString(t *testing.T) {
	tests := []struct {
		DB     *models.DB
		Expect string
	}{
		{
			DB: &models.DB{
				DatabaseUrl:  "postgres://localhost:123",
				DatabaseName: "hello_db",
				UseSSL:       false,
			},
			Expect: "postgres://localhost:123/hello_db?sslmode=disable",
		},
		{
			DB: &models.DB{
				DatabaseUrl:  "postgres://localhost:123",
				DatabaseName: "hello_db",
				UseSSL:       true,
			},
			Expect: "postgres://localhost:123/hello_db",
		},
		{
			DB: &models.DB{
				DatabaseUrl:  "postgres://localhost:123",
				DatabaseName: "",
				UseSSL:       false,
			},
			Expect: "postgres://localhost:123?sslmode=disable",
		},
		{
			DB: &models.DB{
				DatabaseUrl:  "postgres://localhost:123",
				DatabaseName: "",
				UseSSL:       true,
			},
			Expect: "postgres://localhost:123",
		},
		{
			DB: &models.DB{
				DatabaseUrl:  "postgres://localhost:123/",
				DatabaseName: "hello_db",
				UseSSL:       true,
			},
			Expect: "postgres://localhost:123/hello_db",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expect, test.DB.ConnectionString())
	}
}
