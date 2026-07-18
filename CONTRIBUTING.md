# Contributing Guidelines

### note:
This repo was originally cloned from the below repo and deployed to the free tier on heroku 
which ended in 2022. Nearly all references to this have been removed from the code. Some 
references remain as only a reminder of the previous functionality.
https://github.com/heroku/go-getting-started.git

*review this page to smooth out the build/tag/push docker images to the registry
https://devcenter.heroku.com/articles/container-registry-and-runtime#getting-started

## install
docker - to run the api install docker and launch using docker compose

sqlc[https://sqlc.dev/] - install the sqlc binary and extract the binary into the project bin directory

## build

```bash
$ make build
```

## deploy

```bash
$ make deploy
```

## test

```bash
$ make verify
```

token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTc4NDQxMjY2NX0._ze28nnVOdyVFDuiNWJ504cjP0gOd76TihuokRVJMlw;

POST   /account/
    // Create Account - returns initial token=
    curl -X POST -v -d '{"username":"bob","phrase":"13013"}' http://localhost:8080/account/

GET    /account/
    curl -X GET -v --cookie \
    "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTc4NDQxODEwMn0.i4FUmqPrVg4a74fpmuVrXgSYXG39uphSBhnyfYO7-gM" \
      http://localhost:8080/account/

DELETE /account/:id/
    curl -X DELETE -v --cookie \
    "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTc4NDQxODEwMn0.i4FUmqPrVg4a74fpmuVrXgSYXG39uphSBhnyfYO7-gM" \
      http://localhost:8080/account/1/

PATCH  /account/

-----------------------------------------------------------------------------------------------------------------------

POST   /refresh/
POST   /signin/
    // Sign in - returns Set-Token header populated with token=
    curl -v --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTY2NzE1ODQxMn0.d5WA6FOCl_kU4PjP1x0fpsumIPWpSQEn4Fo3MZVuCok" \
    -d '{"username":"bob","phrase":"13013"}' http://localhost:8080/signin/

GET    /welcome/
    // Welcome - can be used if the session cookie is still valid which will issue a new token cookie if needed.

-----------------------------------------------------------------------------------------------------------------------

DELETE /image/
GET    /image/
    // Image - Fetch images list; should check that the jwt is still valid before requesting data using the claimer.
    curl -v -X GET --cookie \ 
      "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTY2NzE1ODQxMn0.d5WA6FOCl_kU4PjP1x0fpsumIPWpSQEn4Fo3MZVuCok" \
      http://localhost:8080/image/

PATCH  /image/
POST   /image/
    // Image - Create image
    curl -v --cookie \ 
     "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImV4cCI6MTY2NzE1ODQxMn0.d5WA6FOCl_kU4PjP1x0fpsumIPWpSQEn4Fo3MZVuCok" \
     -d '{"data":"somefoo"}' http://localhost:8080/image/

