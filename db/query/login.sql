-- -- CREATE TABLE "login" (
-- --   "id"        BIGSERIAL PRIMARY KEY,
-- --   "ipaddress" TEXT NOT NULL DEFAULT '',
-- --   "headers"   TEXT NOT NULL DEFAULT '',
-- --   "token"     VARCHAR NOT NULL,
-- --   "valid"     BOOLEAN NOT NULL DEFAULT TRUE  
-- --   "created"   VARCHAR NOT NULL DEFAULT NOW()
-- -- );

-- -- name: LoginExists :one
-- -- SELECT EXISTS(SELECT 1 FROM account WHERE id = $1);

-- -- name: GetLogin :one
-- SELECT * FROM "login"
-- WHERE id = $1 LIMIT 1;

-- -- name: GetCountLoginByIpAddress :one
-- -- SELECT * FROM "login"
-- -- WHERE ipaddress = $1 LIMIT 1;

-- -- name: ListAccounts :many
-- SELECT * FROM account LIMIT $1 OFFSET $2;

-- -- name: CreateAccount :one
-- INSERT INTO account (
--   username, phrase, salt, created
-- ) VALUES (
--   $1, $2, $3, NOW()
-- )
-- RETURNING *;

-- -- name: SoftDeleteAccount :exec
-- UPDATE account
-- SET deleted = 1
-- WHERE id = $1;

-- -- name: DeleteAccount :exec
-- DELETE FROM account
-- WHERE id = $1;

-- -- name: UpdateAccount :exec
-- UPDATE account
-- SET phrase = $1, updated = NOW()
-- WHERE id = $2
-- RETURNING updated;


