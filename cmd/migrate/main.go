package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ab91dev/codeolx/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("usage: migrate <up | down>")
	}

	cfg := config.MustLoad()

	m, err := migrate.New(
        "file://migrations",
        cfg.DatabaseURL,
	)

	if err != nil {
		log.Fatalf("migration.new: %v",err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	case "down":
		if err := m.Steps(-1); err != nil {
		log.Fatal(err)
	}
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}

	fmt.Println("migration completed successfully")
}
