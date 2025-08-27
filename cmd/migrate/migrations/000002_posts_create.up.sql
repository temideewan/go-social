CREATE TABLE IF NOT EXISTS posts(
id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
title text NOT NULL,
user_id bigint NOT NULL,
content text NOT NULL,
created_at timestamp(0) WITH time zone NOT NULL DEFAULT NOW()
)
