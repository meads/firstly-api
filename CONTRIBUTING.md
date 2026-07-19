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
# database connection string
DATABASE_URL=postgresql://username:password@db:5432/database?sslmode=disable
# secret used for signing jwt tokens
SECRET=

# database initialization values
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=

# username to login on dockerhub 
DOCKER_USERNAME=
```

## build

```bash
$ docker compose build
```

## deploy

```bash
$ make deploy
```

## test

```bash
$ make verify
```

```bash

# POST   /account/
# Create Account - returns initial token=
curl -X POST -v -d '{"username":"bob","phrase":"13013"}' http://localhost:8080/account/

# GET    /account/
curl -X GET -v --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTc4NDQxOTY0NX0.VLD006AUNa_4x-OmqdXECZb6J1yAuouw3JmO1dOwjqw" \
      http://localhost:8080/account/

# DELETE /account/:id/
curl -X DELETE -v --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTc4NDQxOTY0NX0.VLD006AUNa_4x-OmqdXECZb6J1yAuouw3JmO1dOwjqw" \
      http://localhost:8080/account/1/

# PATCH  /account/

# -----------------------------------------------------------------------------------------------------------------------

# POST   /refresh/
# POST   /signin/
# Sign in - returns Set-Token header populated with token=
curl -v --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTY2NzE1ODQxMn0.d5WA6FOCl_kU4PjP1x0fpsumIPWpSQEn4Fo3MZVuCok" \
    -d '{"username":"bob","phrase":"13013"}' http://localhost:8080/signin/

# GET    /welcome/
# Welcome - can be used if the session cookie is still valid which will issue a new token cookie if needed.

# -----------------------------------------------------------------------------------------------------------------------

# DELETE /image/
# GET    /image/
# Image - Fetch images list; should check that the jwt is still valid before requesting data using the claimer.
curl -v -X GET --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTY2NzE1ODQxMn0.d5WA6FOCl_kU4PjP1x0fpsumIPWpSQEn4Fo3MZVuCok" \
      http://localhost:8080/image/

# PATCH  /image/
# POST   /image/
# Image - Create image
curl -v --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTY2NzE1ODQxMn0.d5WA6FOCl_kU4PjP1x0fpsumIPWpSQEn4Fo3MZVuCok" \
     -d '{"data":"somefoo"}' http://localhost:8080/image/
```
