# GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# DOCKER_BUILD=$(shell pwd)/.docker_build
# DOCKER_CMD=$(DOCKER_BUILD)/firstly-api
DATABASE_URL := $(shell grep -iR '^DATABASE_URL*' .env | cut -d= -f2-)
DOCKER_USERNAME := $(shell grep -iR '^DOCKER_USERNAME*' .env | cut -d= -f2-)

# clean:
# 	@rm -rf $(DOCKER_BUILD)
# 	@mkdir -p $(DOCKER_BUILD)

# build: clean
# 	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

test:
	@go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/g'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/g''

test-coverage:
	@go test -v -coverprofile cover.out ./...
	@go tool cover -html=cover.out

# login:
# 	@heroku login
# 	@heroku container:login

# login-docker:
# 	@docker login --username=$(DOCKER_USERNAME) --password=$$(heroku auth:token) registry.heroku.com

local-db-shell:
	@docker exec -it firstly-api-db-1 /bin/bash

local-db-psql:
	@docker exec -it firstly-api-db-1 psql $(DATABASE_URL)

# scale-zero:
# 	@heroku ps:scale api=0

sqlc:
	@sqlc version
	@sqlc compile
	@sqlc generate

tidy:
	@go mod tidy

mockgen:
	@mockgen -package db -destination ./db/querier_mock.go github.com/meads/firstly-api/db Querier
	@mockgen -package security -destination ./security/hmac_mock.go github.com/meads/firstly-api/security Hasher
	@mockgen -package security -destination ./security/claims_mock.go github.com/meads/firstly-api/security Claimer

verify: tidy sqlc mockgen test

migrate-drop-recreate:
	@migrate -path db/migration -database $(DATABASE_URL) drop
	@migrate -path db/migration -database $(DATABASE_URL) up

migrate-drop-recreate-local:
	@migrate -path db/migration -database  $(DATABASE_URL) drop
	@migrate -path db/migration -database  $(DATABASE_URL) up

deploy:
	@git push origin main
# 	@heroku container:push api
# 	@heroku container:release api
