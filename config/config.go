package config

import "github.com/joeshaw/envdecode"

type Config struct {
	ApiUrl        string `env:"MSPL_API_URL,default=http://msplapi.ngrok.io"`
	DatabaseUrl   string `env:"MSPL_DATABASE_URL,default=postgres://localhost:5432/marsupi?sslmode=disable"`
	RabbitMQUrl   string `env:"RABBITMQ_URL,default=amqp://guest:guest@localhost:5672/"`
	WebhookSecret string `env:"WEBHOOK_SECRET,default=8005bcfd-0dda-435f-875b-217929250301"`
}

func NewConfigFromEnvironment() (*Config, error) {
	var cfg Config
	err := envdecode.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
