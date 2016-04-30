all: build

build:
	go-bindata -o migrations.go db/migrations/
	go build
