DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
   id bigserial PRIMARY KEY,
   created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
   email TEXT UNIQUE NOT NULL,
   name TEXT NOT NULL,
   password_hash bytea NOT NULL,
   activated bool NOT NULL DEFAULT false
);

DROP TABLE IF EXISTS posts CASCADE;
CREATE TABLE posts (
   id bigserial PRIMARY KEY,
   created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
   title text UNIQUE NOT NULL,
   url text NOT NULL,
   user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE
);

DROP TABLE IF EXISTS comments CASCADE;
CREATE TABLE comments (
      id bigserial PRIMARY KEY,
      created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
      body text NOT NULL,
      post_id bigint NOT NULL REFERENCES posts ON DELETE CASCADE,
      user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE
);

DROP TABLE IF EXISTS votes CASCADE;
CREATE TABLE votes (
   created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
   user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
   post_id bigint NOT NULL REFERENCES posts ON DELETE CASCADE,
   PRIMARY KEY (user_id, post_id)
);

DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions (
token TEXT PRIMARY KEY,
data BYTEA NOT NULL,
expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

