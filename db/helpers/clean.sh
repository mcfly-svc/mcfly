#!/bin/sh

# delete all builds
psql -d marsupi_test -c 'DELETE FROM build'

# delete all user_project relationships
psql -d marsupi_test -c 'DELETE FROM user_project'

# delete all projects
psql -d marsupi_test -c 'DELETE FROM project'

# delete all users
psql -d marsupi_test -c 'DELETE FROM marsupi_user'