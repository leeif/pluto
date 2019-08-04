package main

import (
	"fmt"
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
	kiper := kiper.NewKiper(filepath.Base(os.Args[0]), "Mercury server")
	kiper.GetKingpinInstance().HelpFlag.Short('h')

	// Init config file from command line and config file
	c := config.GetConfig()

	kiper.SetCommandLineFlag(c, os.Args[1:])

	if err := kiper.ParseCommandLine(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	kiper.SetConfigFilePath(*c.ConfigFile)
	kiper.MergeConfigFile(c)

	if err := rsa.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := database.GetDatabase()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// DB Migration
	if err := migrate.Migrate(db); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start server
	s := server.Server{}
	if err := s.RunServer(); err != nil {
		fmt.Println(err)
	}
}
