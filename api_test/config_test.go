package api_test

import "github.com/mikec/msplapi/config"

func GetTestConfig() *config.Config {
	return &config.Config{
		ApiUrl:      "http://msplapi.ngrok.io",
		DatabaseUrl: "postgres://localhost:5432/marsupi_test?sslmode=disable",
	}
}
