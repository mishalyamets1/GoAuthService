CREATE SCHEMA sso;
CREATE TABLE sso.users (
    id        BIGSERIAL PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    is_admin  BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE TABLE sso.apps (
    id     SERIAL PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL
);