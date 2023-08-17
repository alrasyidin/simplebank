# Simple Bank Rest API
This Rest API built with Golang And SQLC, User can have many accounts with different currencies. User can perform transfer from two accounts.

## Features

* Create and mange bank accounts.
* Records history balance changes to each of accounts.
* Perform transfer money from two account.
* Send email with background workers.
* GRPC APIs
* Perform Unit and Integration testing

## Technologies

* Go
* PostgreSQL (generated with SQLC)
* asynq for sending email or background process
* GRPC API (for several endpoints)
* JWT
* Swagger
* Docker

## How to running

### Using your machine
- install air for hot reloading golang apps [(air comstrek)](https://github.com/cosmtrek/air)
- install PostgreSQL on your machine
- install SQLC binary [(sqlc)](https://docs.sqlc.dev/en/stable/overview/install.html)
- perform migration ```make migrateup```
- generate sqlc ```make sqlc```
- generate mock for testing ```make mock```
- running apps with ```make server```

### Using Docker
- install Docker
- install air for hot reloading golang apps [(air comstrek)](https://github.com/cosmtrek/air)
- install SQLC binary [(sqlc)](https://docs.sqlc.dev/en/stable/overview/install.html
- running docker compose with ```docker compose up -d```
- perform migration ```make migrateup```
- generate sqlc ```make sqlc```
- generate mock for testing ```make mock```
- running apps with ```make server```

