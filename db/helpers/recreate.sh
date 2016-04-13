#!/bin/sh

# cd to the path of this script: marsupi-api/db/helpers/recreate.sh
parent_path=$( cd "$(dirname "${BASH_SOURCE}")" ; pwd -P )
cd "$parent_path"

# delete database schema and data by running all migrations down
migrate -url "postgres://localhost:5432/marsupi_test?sslmode=disable" -path ../migrations down

# recreate the database by running migrations up
migrate -url "postgres://localhost:5432/marsupi_test?sslmode=disable" -path ../migrations up