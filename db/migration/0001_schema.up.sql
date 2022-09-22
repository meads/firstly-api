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
  "phrase"   TEXT       NOT NULL,
  "salt"     TEXT       NOT NULL,
  "created"  VARCHAR NOT NULL,
  "updated"  VARCHAR NOT NULL,
  "deleted"  BOOLEAN NOT NULL DEFAULT FALSE
);
