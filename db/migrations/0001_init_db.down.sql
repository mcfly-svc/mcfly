DROP INDEX IF EXISTS project_service_username_name;
DROP INDEX IF EXISTS provideraccesstoken_user;
DROP INDEX IF EXISTS project_handle_provider;

DROP TABLE IF EXISTS build;
DROP TABLE IF EXISTS user_project;
DROP TABLE IF EXISTS project;
DROP TABLE IF EXISTS provider_access_token;
DROP TABLE IF EXISTS marsupi_user;

DROP TYPE IF EXISTS provider;
DROP TYPE IF EXISTS build_deploy_status;
