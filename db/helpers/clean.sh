#!/bin/sh

# delete all projects
psql -d marsupi_test -c 'DELETE FROM project'