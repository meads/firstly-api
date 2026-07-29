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
# api database connection string
DATABASE_URL=postgresql://username:password@db:5432/databasename?sslmode=disable

# used for jwt signing
SECRET=

# used for configuring gin router cors middleware
ALLOW_ORIGINS=

# variables for db service
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=


```

## generate

Go module housekeeping is performed. Generates the data access code from the 
db project using sqlc and generates the mocks used in tests for the data access code using 
mockgen.

```bash
$ make generate
```

To generate new dml sql scripts use golang-migrate/migrate tool. This will add files to the
db/migration directory following the migrate tools naming convention. 

```bash
$ migrate create -ext sql -dir db/migration create_example_table
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

# POST   /register/
# Create User - returns initial token in Authorization header
curl -X POST -v -d '{"username":"bob","password":"13013"}' http://localhost:8080/register/

# Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJib2IiLCJ0eXBlIjoiYWNjZXNzIiwic3ViIjoiYm9iIiwiZXhwIjoxNzg1MTgzNjc1LCJpYXQiOjE3ODUxODMzNzUsImp0aSI6IjdjZmJiNTU2LTcyYzktNGMzYy04NWQ2LWRhNjkyNzdhODZjNyJ9.qmHS78GyVOXIR-WJ546qqOzOXosBLbK8Bi2O7q_YLQM

# Create session - returns sessionId, accessToken, refreshToken, accessTokenExpiresAt, refreshTokenExpiresAt, username in body
curl -X POST -v -d '{"username":"bob","password":"13013"}' http://localhost:8080/login/
# {"sessionId":"4346a042-5dde-468b-8ad1-9b10c4d71d8e","accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJib2IiLCJ0eXBlIjoiYWNjZXNzIiwic3ViIjoiYm9iIiwiZXhwIjoxNzg1MTgzNzI3LCJpYXQiOjE3ODUxODM0MjcsImp0aSI6IjBkYzEyODg5LThiNTAtNDQzMi04ZGY0LTZkMWM1N2I5MDVhYyJ9.g6oKiiLayHDS_uVzKMJB7cvPrO702osiFjk1yVwZBfE","refreshToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJib2IiLCJ0eXBlIjoicmVmcmVzaCIsInN1YiI6ImJvYiIsImV4cCI6MTc4NTI2OTgyNywiaWF0IjoxNzg1MTgzNDI3LCJqdGkiOiI0MzQ2YTA0Mi01ZGRlLTQ2OGItOGFkMS05YjEwYzRkNzFkOGUifQ.QQnlW3Xq5OzZk5LvBQeqU3IMJuVn6lc8ldbMeJQf_W4","accessTokenExpiresAt":"2026-07-27T20:22:07Z","refreshTokenExpiresAt":"2026-07-28T20:17:07Z","username":"bob"}

# Create new access token from the supplied refreshToken - returns accessToken, accessTokenExpiresAt in body
curl -X POST -v -d '{"refreshToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJib2IiLCJ0eXBlIjoicmVmcmVzaCIsInN1YiI6ImJvYiIsImV4cCI6MTc4NTI3Mzg0OSwiaWF0IjoxNzg1MTg3NDQ5LCJqdGkiOiI1ZTY3YWJmMC0zYzkzLTQ1YjUtOWJjNS1iODUxMWM4OGUwMTUifQ.lDUod1IhID-nBP3qxrX8mlGAS85cAqIdeElPhvkUYJA"}' http://localhost:8080/refresh/

# {"accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJib2IiLCJzdWIiOiJib2IiLCJleHAiOjE3ODUwOTI3ODQsImlhdCI6MTc4NTA5MTg4NCwianRpIjoiMGFiNTkzNGUtZmFjYi00NDY4LTgzOGQtZDY1NThiZDUzYmFmIn0.v7UdI4_Z9BdQqCU_YxSEsaCF2jGGBwBnu4c1As_GYoc","accessTokenExpiresAt":"2026-07-26T19:06:24Z"}

curl -X POST -v http://localhost:8080/logout/6f5f326f-63cf-4b7f-9891-21dd3a516b76

curl -X POST -v http://localhost:8080/revoke/4346a042-5dde-468b-8ad1-9b10c4d71d8e

# GET    /users/
curl -X GET -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJib2IiLCJzdWIiOiJib2IiLCJleHAiOjE3ODUwOTQ0NDAsImlhdCI6MTc4NTA5MzU0MCwianRpIjoiZDc0NDc5NzgtMTM2OC00M2ZhLThhNzUtMDcxMDgwMDQ4OTQ3In0.qpRYE7glGnW9fZlI6kvmgm2VVZOlWugD5jEIklCVt1c" \
      http://localhost:8080/users/

# DELETE /users/:id/
curl -X DELETE -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJib2IiLCJzdWIiOiJib2IiLCJleHAiOjE3ODUwOTM0ODksImlhdCI6MTc4NTA5MzE4OSwianRpIjoiYjYzYzJmNzYtYjcyMy00ZjhmLWJjYjgtZWUwMTMxODhlYTU2In0.VnQSQuJpAI9u8qgADmYArgzKF6AnqIeaQ8uA_8BtC7A" \
      http://localhost:8080/users/1/

# PATCH  /users/
# -----------------------------------------------------------------------------------------------------------------------

# GET /protected/
# return api JSON data from a route that is protected by JWT validation middleware
curl -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJib2IiLCJ0eXBlIjoiYWNjZXNzIiwic3ViIjoiYm9iIiwiZXhwIjoxNzg1MTg4NDA0LCJpYXQiOjE3ODUxODc1MDQsImp0aSI6IjUwM2JjYTgwLThmYzItNDc4OS04MGY4LTVhN2NjMmY5ZjU1OSJ9.M3XdWfvaYqmkAtEPbNgFbBFpojOA7bcZFZgX_GCH-AI" -v http://localhost:8080/protected/

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
