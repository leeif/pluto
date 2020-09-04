package main

import (
	"log"
	"os"

	"github.com/MuShare/pluto/config"

	"github.com/MuShare/pluto/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/MuShare/pluto/utils/migrate"
)

func main() {
	cfg, err := config.NewConfig(os.Args, "")

	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDatabase(cfg)

	if err != nil {
		log.Fatal(err)
	}

	if err := migrate.Migrate(db); err != nil {
		log.Fatal(err)
	}

	log.Println("migrate success")
}
