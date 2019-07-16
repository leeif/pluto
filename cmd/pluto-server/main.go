package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/server"

	"github.com/leeif/pluto/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "Mercury server")
	a.Version("0.0.1")
	a.HelpFlag.Short('h')

	c := config.NewConfig()
	c.Parse(a, os.Args[1:])

	logger := log.NewLogger(c.Log)

	s := server.NewServer(c.Server, logger)
	if err := s.RunServer(); err != nil {
		fmt.Println(err)
	}
}
