#!/bin/sh

psql -d marsupi_test <<- EOF
	DELETE FROM build;
	DELETE FROM user_project;
	DELETE FROM project;
	DELETE FROM provider_access_token;
	DELETE FROM marsupi_user;
EOF
