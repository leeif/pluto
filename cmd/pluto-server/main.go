package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/leeif/kiper"
	"github.com/leeif/pluto/database"
	"github.com/leeif/pluto/utils/migrate"

	"github.com/leeif/pluto/config"

	"github.com/leeif/pluto/server"

	"github.com/leeif/pluto/utils/rsa"
)

func main() {
	kiper := kiper.NewKiper(filepath.Base(os.Args[0]), "Pluto server")
	kiper.GetKingpinInstance().HelpFlag.Short('h')

	// Init config file from command line and config file
	c := config.GetConfig()

	if err := kiper.ParseCommandLine(c, os.Args[1:]); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := kiper.ParseConfigFile(*c.ConfigFile); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := kiper.MergeConfigFile(c); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := config.CheckConfig(c); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := rsa.Init(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	db, err := database.GetDatabase()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// DB Migration
	if err := migrate.Migrate(db); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Start server
	s := server.Server{}
	err = s.RunServer()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
