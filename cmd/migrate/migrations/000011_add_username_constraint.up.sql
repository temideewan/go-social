ALTER TABLE users ADD CONSTRAINT users_username_key unique(username);
ALTER TABLE users ALTER COLUMN username SET NOT NULL;
