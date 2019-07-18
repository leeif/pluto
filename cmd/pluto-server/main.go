package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leeif/pluto/config"

	"github.com/leeif/pluto/server"

	"github.com/leeif/pluto/utils/rsa"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "Mercury server")
	a.Version("0.0.1")
	a.HelpFlag.Short('h')

	// Init config file from command line and config file
	c := config.GetConfig()
	c.Parse(a, os.Args[1:])

	if err := rsa.Init(); err != nil {
		fmt.Println(err)
	}

	s := server.NewServer()
	if err := s.RunServer(); err != nil {
		fmt.Println(err)
	}
}
