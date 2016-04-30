#!/bin/sh

echo "cleaning the Database"
marsupi-api db.clean

echo "change supported providers to test values"
psql -d marsupi_test <<- EOF
	ALTER TYPE provider ADD VALUE 'jabroni.com';
	ALTER TYPE provider ADD VALUE 'schlockbox';
EOF

echo "adding seed users"

psql -d marsupi_test <<- EOF

	INSERT INTO marsupi_user (name, access_token) 
	VALUES ('Matt Mockerson', 'mock_token_123')
	RETURNING id;
	\gset
	\echo new user id=:id

	INSERT INTO provider_access_token (provider, provider_username, access_token, user_id)
	VALUES ('jabroni.com', 'mattmocks', 'mock_jabroni.com_token_123', :id);
	\echo new provider_access_token 'mock_jabroni.com_token_123'

	INSERT INTO provider_access_token (provider, provider_username, access_token, user_id)
	VALUES ('schlockbox', 'mattmocks@gmail.com', 'mock_schlockbox_token_123', :id);
	\echo new provider_access_token 'mock_schlockbox_token_123'

EOF
