package models_test

import "github.com/mcfly-svc/mcfly/config"

func GetTestConfig() *config.Config {
	return &config.Config{
		ApiUrl:         "http://mcfly.ngrok.io",
		DatabaseUrl:    "postgres://localhost:5432",
		DatabaseName:   "mcfly_test",
		DatabaseUseSSL: false,
		WebhookSecret:  "mock_webhook_secret",
	}
}
