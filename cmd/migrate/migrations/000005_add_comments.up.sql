CREATE TABLE
  IF NOT EXISTS comments (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    post_id bigint NOT NULL,
    user_id bigint NOT NULL,
    content TEXT NOT NULL,
    created_at timestamp(0)
    WITH
      time zone NOT NULL DEFAULT NOW ()
  );
