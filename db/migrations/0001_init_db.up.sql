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
 	user_id									integer REFERENCES marsupi_user(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX provideraccesstoken_user ON provider_access_token (provider, user_id);

CREATE TABLE project (
	id											serial PRIMARY KEY,
	handle									text,
	source_url							text,
	source_provider					provider
);
CREATE UNIQUE INDEX project_handle_provider ON project (handle, source_provider);

CREATE TABLE build (
	id											serial PRIMARY KEY,
	handle 									text,
	project_id							integer REFERENCES project(id) ON DELETE CASCADE
);

CREATE TABLE user_project (
	user_id									integer REFERENCES marsupi_user(id) ON DELETE CASCADE,
	project_id							integer REFERENCES project(id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, project_id)
);
