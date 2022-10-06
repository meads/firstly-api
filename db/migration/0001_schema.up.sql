CREATE TABLE "image" (
  "id"   BIGSERIAL  PRIMARY KEY,
  "data" TEXT       NOT NULL,
  "memo" TEXT       NOT NULL,
  "created" VARCHAR NOT NULL,
  "updated" VARCHAR NOT NULL,
  "deleted" BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE "account" (
  "id"       BIGSERIAL  PRIMARY KEY,
  "username" TEXT       NOT NULL,
  "phrase"   BYTEA      NOT NULL,
  "salt"     TEXT       NOT NULL,
  "created"  VARCHAR NOT NULL,
  "updated"  VARCHAR NOT NULL DEFAULT '',
  "deleted"  BOOLEAN NOT NULL DEFAULT FALSE
);

-- CREATE TABLE "login" (
--   "id"        BIGSERIAL PRIMARY KEY,
--   "ipaddress" TEXT NOT NULL DEFAULT '',
--   "headers"   TEXT NOT NULL DEFAULT '',
--   "token"     VARCHAR NOT NULL,
--   "valid"     BOOLEAN NOT NULL DEFAULT TRUE  
--   "created"   VARCHAR NOT NULL DEFAULT NOW()
-- );

