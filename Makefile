all: get build test

get:
	go get -u github.com/jteeuwen/go-bindata/...

build:
	go-bindata -pkg db -o db/assets.go db/migrations/
	go install

database:
	msplapi create-db
	msplapi seed-db

test: database
	go test ./...

run: build
	msplapi run