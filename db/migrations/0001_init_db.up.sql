CREATE TYPE provider AS ENUM ('github', 'dropbox', 'travis', 'jenkins');
CREATE TYPE build_deploy_status AS ENUM('succeeded', 'pending', 'failed');

CREATE TABLE mcfly_user (
	id											serial PRIMARY KEY,
	name										text,
	access_token						text
);

CREATE TABLE provider_access_token (
	provider								provider,
	provider_username 			text,
	access_token 						text,
 	user_id									integer REFERENCES mcfly_user(id) ON DELETE CASCADE
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
	project_id							integer REFERENCES project(id) ON DELETE CASCADE,
	deploy_status 					build_deploy_status,
	provider_url   					text
);

CREATE TABLE user_project (
	user_id									integer REFERENCES mcfly_user(id) ON DELETE CASCADE,
	project_id							integer REFERENCES project(id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, project_id)
);
