all: build test

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