CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users(
id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
email citext UNIQUE NOT NULL,
username VARCHAR(255),
password bytea NOT NULL,
created_at timestamp(0) WITH time zone NOT NULL DEFAULT NOW()
)
-- first_name VARCHAR(255) NOT NULL,
-- last_name VARCHAR(255) NOT NULL,
