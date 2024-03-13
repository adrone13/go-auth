package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
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
}

func down() {
	fmt.Println("Migrating down")
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

	files, err := os.ReadDir("./migrations")
	if err != nil {
		panic(err)
	}

	last := 0

	for _, f := range files {
		parts := strings.Split(f.Name(), "_")

		n, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("Invalid migration number: %s", parts[0]))
		}

		fmt.Println(n, reflect.TypeOf(n))

		if n > last {
			last = n
		}
	}

	upName := fmt.Sprintf("%d_%s.up.sql", last+1, name)
	downName := fmt.Sprintf("%d_%s.down.sql", last+1, name)

	fmt.Printf("New migration name:\n\t%s\n\t%s", upName, downName)

	err = os.WriteFile(fmt.Sprintf("./migrations/%s", upName), make([]byte, 0), 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(fmt.Sprintf("./migrations/%s", downName), make([]byte, 0), 0644)
	if err != nil {
		panic(err)
	}
}
