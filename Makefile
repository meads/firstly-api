GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/firstly-api

clean:
	@rm -rf $(DOCKER_BUILD)
	@mkdir -p $(DOCKER_BUILD)

build: clean
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

test:
	@go test -v ./...

login:
	@heroku login
	@heroku container:login

login-docker:
	@docker login --username=resetheadhat@gmail.com --password=$$(heroku auth:token) registry.heroku.com

local:
	@docker-compose up

local-clean: clean
	@docker-compose build

local-db-shell:
	@docker exec -it firstly-api_db_1 /bin/bash

local-db-psql:
	@docker exec -it firstly-api_db_1 psql postgresql://postgres:password@localhost:5432/postgres

scale-zero:
	@heroku ps:scale web=0

sqlc:
	@sqlc version
	@sqlc compile
	@sqlc generate

tidy:
	@go mod tidy

mockgen:
	@mockgen -package db -destination ./db/store_mock.go github.com/meads/firstly-api/db Store

verify: tidy sqlc mockgen build test

deploy:
	@ssh-session
	@git push origin main
	@heroku container:push web
	@heroku container:release web

