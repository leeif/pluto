package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/leeif/pluto/server"

	plog "github.com/leeif/pluto/log"
	"github.com/leeif/pluto/route"

	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
	"github.com/leeif/pluto/config"
	"go.uber.org/fx"

	"github.com/leeif/pluto/database"
	"github.com/leeif/pluto/utils/migrate"

	"github.com/leeif/pluto/utils/rsa"
	"github.com/leeif/pluto/utils/view"
)

var VERSION = ""

func printConfig(config *config.Config, logger *plog.PlutoLog) {
	logger.Info(fmt.Sprintf("AccessTokenExpire: %d", config.JWT.AccessTokenExpire))
	logger.Info(fmt.Sprintf("RegisterVerifyTokenExpire: %d", config.JWT.RegisterVerifyTokenExpire))
	logger.Info(fmt.Sprintf("ResetPasswordResultTokenExpire: %d", config.JWT.ResetPasswordResultTokenExpire))
	logger.Info(fmt.Sprintf("ResetPasswordTokenExpire: %d", config.JWT.ResetPasswordTokenExpire))
}

func register(router *mux.Router, db *gorm.DB, config *config.Config, logger *plog.PlutoLog) error {

	printConfig(config, logger)

	if err := rsa.Init(config, logger); err != nil {
		logger.Error(err.Error())
		return err
	}

	if err := view.Init(config, logger); err != nil {
		logger.Error(err.Error())
		return err
	}

	// DB Migration
	if err := migrate.Migrate(db); err != nil {
		logger.Error(err.Error())
		return err
	}

	// add router
	route.Router(router, db, config, logger)
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
