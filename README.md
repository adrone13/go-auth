# Golang Auth service
Learning project to get familiar with Go programming language utilizing __DDD__ practices and minimal dependencies.

Service features:
* Sign up/Sign in
* Access token rotation
* Refresh token rotation with replay attack detection
* Simple custom [gojwt](https://github.com/adrone13/gojwt) package

## Project structure
* cmd
  * auth - app entry point
  * migrate - migration script written using [golang-migrate](https://github.com/golang-migrate/migrate) library to make migrations part of the project rather than use external tools.
* internal
  * app - application's business logic
  * db, logger, server, etc. - projects infrustructure and various helpers

## Commands
```bash
# build
$ make build

# run
$ make run
$ make watch # watch mode

# db
$ make db-run         # run db docker container
$ make db-stop        # stop db container
$ make migrate-create # create next migration file pair (up and down)
$ make migrate-up     # run all unapplied migrations
$ make migrate-down   # revert all migrations

# test
make test

# binary clean up
make clean
```

