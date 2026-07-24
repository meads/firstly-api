# Contributing Guidelines

### note:
This repo was originally cloned from the below repo and deployed to the free tier on heroku 
which ended in 2022. Nearly all references to this have been removed from the code. Some 
references remain as only a reminder of the previous functionality.
https://github.com/heroku/go-getting-started.git

*review this page to smooth out the build/tag/push docker images to the registry
https://devcenter.heroku.com/articles/container-registry-and-runtime#getting-started

## install

make sure you have the following in your .bashrc or .bash_profile 

```bash
# Adds the main 'go' compiler command to your path (probably already set)
export PATH="$PATH:/usr/local/go/bin"

# export the GOPATH variable explicitly
export GOPATH="$HOME/go"

# Adds user installed binaries 'go install' tools are in your path
export PATH="$PATH:$GOPATH/bin"
```
docker
  - why: To run the api and database in containers. Install docker desktop.
  - how: [typical docker install steps](https://www.docker.com/get-started/)

sqlc
  - why: to generate data access code for the database
  - how: 

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

mockgen
  - why: to generate mocks of the Querier interface for mocking data calls in http handler unit tests
  - how: 

```bash  
go install go.uber.org/mock/mockgen@latest
```

migrate
  - why: to control the version and state of the database in use
      Dual Files: Generates distinct .up.sql (to apply changes) and .down.sql (to roll back changes) schema files.
      Safety Tracking: Creates a tracking table (schema_migrations) in your database to log the exact current version and flag failed ("dirty") states.
  - how: 

```bash  
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## environment variables

Create a local .env file in the root directory and specify the values below. 
The .env file is excluded from the project via .gitignore file.
    
```env
# variables for api service
DATABASE_URL=postgresql://username:password@db:5432/databasename?sslmode=disable
# used in creation of hmac for generating and validating password hashes
SECRET=
# origins to allow in cors requests using gin middleware comma separated. 
# don't leave trailing commas
ALLOW_ORIGINS=

# variables for db service
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=

# variables for makefile
DOCKER_USERNAME=

```

## generate

Go module housekeeping is performed. Generates the data access code from the 
db project using sqlc and generates the mocks used in tests for the data access code using 
mockgen.

```bash
$ make generate
```

## test

Run the unit tests for the entire application.

```bash
$ make test
```

## test coverage

Run the unit tests for the entire application and generate a coverage report.
```bash
$ make test-cover
```

## verify

Quickly run both generating code and test recipes
```bash
$ make verify
```

## build

Build the docker containers with the appropriate files and environment 
configurations. 

```bash
$ docker compose build
```

## run local

Run local instances of the docker compose containers

```bash
$ docker compose up
```

## cleanup local

Remove all containers and start fresh.
```bash
$ docker compose down
```

## deploy

```bash
# TODO: create deploy flow
$ make deploy
```


```bash

# POST   /account/
# Create Account - returns initial token in Authorization header
curl -X POST -v -d '{"username":"bob","password":"13013"}' http://localhost:8080/account/

# GET    /account/
curl -X GET -v -H "Authorization: Bearer [token here]" \
      http://localhost:8080/account/

# DELETE /account/:id/
curl -X DELETE -v -H "Authorization: Bearer [token here]" \
      http://localhost:8080/account/1/

# PATCH  /account/
# -----------------------------------------------------------------------------------------------------------------------

# GET /protected/
# return api JSON data from a route that is protected by JWT validation middleware
curl -X GET -H "Authorization: Bearer [token here]" -v http://localhost:8080/protected/

# -----------------------------------------------------------------------------------------------------------------------

# POST   /refresh/
# POST   /signin/
# Sign in - returns Set-Token header populated with token=
curl -X POST -v -H "Authorization: Bearer [token here]" \
    -d '{"username":"bob","password":"13013"}' http://localhost:8080/signin/

# GET    /welcome/
# Welcome - can be used if the session cookie is still valid which will issue a new token cookie if needed.

# -----------------------------------------------------------------------------------------------------------------------

# # DELETE /image/
# # GET    /image/
# # Image - Fetch images list; should check that the jwt is still valid before requesting data using the claimer.
# curl -v -X GET  http://localhost:8080/image/

# # PATCH  /image/
# # POST   /image/
# # Image - Create image
# curl -v -d '{"data":"somefoo"}' http://localhost:8080/image/
# ```
