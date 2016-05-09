CREATE TYPE provider AS ENUM ('github', 'dropbox', 'travis', 'jenkins');

CREATE TABLE marsupi_user (
	id											serial PRIMARY KEY,
	name										text,
	access_token						text
);

CREATE TABLE provider_access_token (
	provider								provider,
	provider_username 			text,
	access_token 						text,
 	user_id									integer REFERENCES marsupi_user(id)
);
CREATE UNIQUE INDEX provideraccesstoken_user ON provider_access_token (provider, user_id);

CREATE TABLE project (
	id											serial PRIMARY KEY,
	handle									text,
	source_url							text,
	source_provider					provider
);

CREATE TABLE build (
	id											serial PRIMARY KEY,
	project_id							integer REFERENCES project(id),
	hash										text,
	build_provider					provider
);

CREATE TABLE user_project (
	user_id									integer REFERENCES marsupi_user(id),
	project_id							integer REFERENCES project(id),
	PRIMARY KEY (user_id, project_id)
);
