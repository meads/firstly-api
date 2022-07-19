GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/firstly-api

# Removing this target may have side effects on the ci build/deploy pipeline.
$(DOCKER_CMD): clean
	@mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

clean:
	@rm -rf $(DOCKER_BUILD)

build: clean
	@mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .


test: sqlc mockgen
	@go test -v

release:
	@heroku container:push web

# Removing this target may have side effects on the ci build/deploy pipeline.
heroku: $(DOCKER_CMD)
	@heroku container:push web

local: build
	@heroku local web

scale-zero:
	@heroku ps:scale web=0

sqlc:
	@sqlc version
	@sqlc compile
	@sqlc generate

mockgen:
	@mockgen -package mockapi -destination ./db/mock/store.go github.com/meads/firstly-api/db/api Store

