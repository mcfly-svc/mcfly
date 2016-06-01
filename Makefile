all: build test

get:
	go get -u github.com/stretchr/testify/mock
	go get -u github.com/jteeuwen/go-bindata/...

build:
	go-bindata -pkg db -o db/assets.go db/migrations/
	go install

test: build
	mcfly init-db mcfly_test
	mcfly seed-db mcfly_test
	go test ./...

run: build
	mcfly run

database:
	mcfly create-db mcfly
	mcfly init-db mcfly
	mcfly seed-db mcfly
	mcfly create-db mcfly_test
	mcfly init-db mcfly_test
	mcfly seed-db mcfly_test

setup: get build database

run-ngrok:
	ngrok http -subdomain=mcfly 8081

# gofmt
# govet