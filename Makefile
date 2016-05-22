all: get build test

get:
	go get -u github.com/stretchr/testify/mock
	go get -u github.com/jteeuwen/go-bindata/...

build:
	go-bindata -pkg db -o db/assets.go db/migrations/
	go install

test: build
	msplapi init-db marsupi_test
	msplapi seed-db marsupi_test
	go test ./...

run: build
	msplapi run

database:
	msplapi create-db marsupi
	msplapi init-db marsupi
	msplapi seed-db marsupi
	msplapi create-db marsupi_test
	msplapi init-db marsupi_test
	msplapi seed-db marsupi_test

ngrok:
	ngrok http -subdomain=msplapi 8081

# gofmt
# govet