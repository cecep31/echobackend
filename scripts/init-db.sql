-- Executed once by postgres on first container start (docker-entrypoint-initdb.d).
-- Goose stores its version history in the non-default `custom` schema
-- (GOOSE_TABLE=custom.goose_migrations), so the schema must exist before
-- the first `goose up`.
CREATE SCHEMA IF NOT EXISTS custom;
