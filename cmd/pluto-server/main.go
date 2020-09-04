package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/MuShare/pluto/manage"
	"github.com/MuShare/pluto/utils/admin"

	"github.com/MuShare/pluto/server"

	plog "github.com/MuShare/pluto/log"
	"github.com/MuShare/pluto/route"

	"github.com/MuShare/pluto/config"
	"go.uber.org/fx"

	"github.com/MuShare/pluto/database"

	_ "github.com/go-sql-driver/mysql"
	perror "github.com/MuShare/pluto/datatype/pluto_error"
	"github.com/MuShare/pluto/utils/rsa"
	"github.com/MuShare/pluto/utils/view"
)

// VERSION is the pluto build version
var VERSION = ""

func register(router *route.Router, db *sql.DB, config *config.Config) error {

	if err := rsa.Init(config); err != nil {
		log.Fatalln(err.Error())
		return err
	}

	if err := view.InitView(config); err != nil {
		log.Fatalln(err.Error())
		return err
	}

	if err := admin.Init(db, config); err != nil {
		if err.PlutoCode == perror.ServerError.PlutoCode {
			log.Fatalln(err.LogError.Error())
			return err.LogError
		}
	}

	// register routes
	router.Register()

	return nil
}

func main() {

	app := fx.New(
		fx.Provide(
			func() []string {
				return os.Args
			},
			func() string {
				return VERSION
			},
			config.NewConfig,
			database.NewDatabase,
			plog.NewLogger,
			server.NewMux,
			route.NewRouter,
			manage.NewManager,
		),
		fx.Invoke(register),
		fx.NopLogger,
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
