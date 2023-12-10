# Auth service written in Go lang
The main goal of this project is to learn Go by making a service with minimal dependencies (`joho/godotenv` and `x/crypto` being the only two at the moment) and eventually develop it in full microservice driven app.

## Structure
Go is always described to be a simple language and that any extra abstraction should be omitted, so in this project I've decided to keep abstractions minimal, for instance, using interfaces for repositories. Since I enjoy a lot of __DDD__ practices I decided to have an __Application__ layer. It is represented by `app` directory and defines entities (like `user.go`), interfaces (like `user_repository.go`) and use cases (like `log_in_usecase.go`).

Other directories mostly serve infrastructure functions like HTTP server, DB connection, JWT signing/parsing, etc.

## Commands
```bash
# build
$ make
$ make all
$ make build

# run
$ make run

# test
make test

# binary clean up
make clean
```

