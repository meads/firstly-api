CREATE TABLE "images" (
  "id"   BIGSERIAL PRIMARY KEY,
  "data" TEXT      NOT NULL,
  "created" VARCHAR(8) NOT NULL,
  "deleted" INT DEFAULT 0
);
