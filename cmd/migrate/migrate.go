package main

import (
	"auth/internal/config"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var url = fmt.Sprintf(
	"postgres://%s:%s@%s:%d/%s?sslmode=disable",
	config.Values.DbUser,
	config.Values.DbPassword,
	config.Values.DbHost,
	config.Values.DbPort,
	config.Values.DbName,
)

func main() {
	args := os.Args[1:]

	fmt.Println("Args:", args)

	if len(args) == 0 {
		log.Fatalln("Command required")
	}

	command := args[0]

	switch command {
	case "create":
		create()
	case "up":
		up()
	case "down":
		down()
	default:
		log.Fatalln("Invalid command")
	}
}

func up() {
	fmt.Println("Migrating up")

	migration, err := migrate.New("file://cmd/migrate/migrations", url)
	if err != nil {
		log.Fatal(err)
	}

	err = migration.Up()
	if err != nil {
		log.Fatal(err)
	}
}

func down() {
	fmt.Println("Migrating down")

	migration, err := migrate.New("file://cmd/migrate/migrations", url)
	if err != nil {
		log.Fatal(err)
	}

	err = migration.Down()
	if err != nil {
		log.Fatal(err)
	}
}

func create() {
	if len(os.Args) < 3 {
		panic("Invalid args")
	}

	name := os.Args[2]
	if name == "" {
		panic("Invalid migration name")
	}

	fmt.Println("Creating migration...")

	files, err := os.ReadDir("./cmd/migrate/migrations")
	if err != nil {
		panic(err)
	}

	last := 0

	for _, f := range files {
		parts := strings.Split(f.Name(), "_")

		n, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatalf("Invalid migration number: %s\n", parts[0])
		}

		if n > last {
			last = n
		}
	}

	upName := fmt.Sprintf("%d_%s.up.sql", last+1, name)
	downName := fmt.Sprintf("%d_%s.down.sql", last+1, name)

	fmt.Printf("Migration created:\n\t%s\n\t%s\n", upName, downName)

	err = os.WriteFile(fmt.Sprintf("./cmd/migrate/migrations/%s", upName), make([]byte, 0), 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(fmt.Sprintf("./cmd/migrate/migrations/%s", downName), make([]byte, 0), 0644)
	if err != nil {
		panic(err)
	}
}
