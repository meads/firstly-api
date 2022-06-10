CREATE TABLE images (
  "id"   BIGSERIAL PRIMARY KEY,
  "name" TEXT      NOT NULL,
  "data" TEXT NOT NULL,
  "created" NVARCHAR(8) NOT NULL,
  "deleted" BOOLEAN DEFAULT FALSE
);
-- 	host     = "ec2-52-204-195-41.compute-1.amazonaws.com"
-- 	port     = 5432
-- 	user     = "irdanpwkdvbzxg"
-- 	password = "c29b40a6619957f7b572795b81f7805414be54ccd3675c07761f4cc61894d83e"
-- 	dbname   = "d89dudhb3lei05"
