all: get build database test

get:
	go get -u github.com/jteeuwen/go-bindata/...

build:
	go-bindata -pkg db -o db/assets.go db/migrations/
	go install

database:
	msplapi create-db
	msplapi seed-db

test: build
	go test ./...

run: build
	msplapi run

ngrok:
	ngrok http -subdomain=msplapi 8081