marsupi-api
-----------

Data access layer for marsupi


Setup
=====

go get -u github.com/jteeuwen/go-bindata/...


Run API Server
==============

marsupi-api run


Testing
=======

go test ./...


Migrations
==========

go get -u github.com/mattes/migrate

migrate -url postgres://localhost:5432/marsupi?sslmode=disable -path ./db/migrations create add_field_to_table
migrate -url postgres://localhost:5432/marsupi?sslmode=disable -path ./db/migrations up
migrate -url postgres://localhost:5432/marsupi_test?sslmode=disable -path ./db/migrations up