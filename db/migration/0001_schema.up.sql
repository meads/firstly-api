-- CREATE TABLE "image" (
--   "id"   BIGSERIAL  PRIMARY KEY,
--   "data" TEXT       NOT NULL,
--   "memo" TEXT       NOT NULL DEFAULT '',
--   "created" VARCHAR NOT NULL,
--   "updated" VARCHAR NOT NULL DEFAULT '',
--   "deleted" BOOLEAN NOT NULL DEFAULT FALSE
-- );

-- CREATE TABLE users (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     email VARCHAR(255) UNIQUE NOT NULL,
--     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

CREATE TABLE account (
  id       BIGSERIAL  PRIMARY KEY,
  username TEXT       NOT NULL,
  "password"   BYTEA      NOT NULL,
  salt     TEXT       NOT NULL,
  created  VARCHAR NOT NULL,
  updated  VARCHAR NOT NULL DEFAULT '',
  deleted  BOOLEAN NOT NULL DEFAULT FALSE
);

-- TODO implement complete jwt access/refresh token authentication
-- CREATE TABLE sessions (
--     id  BIGSERIAL PRIMARY KEY,
--     account_id BIGSERIAL NOT NULL REFERENCES account(id) ON DELETE CASCADE,
--     token_hash TEXT NOT NULL,
--     expires_at VARCHAR NOT NULL,
--     created_at VARCHAR NOT NULL DEFAULT NOW()
-- );

-- -- Index the foreign key for faster lookups when querying an account's sessions
-- CREATE INDEX idx_sessions_account_id ON jwt_sessions(account_id);