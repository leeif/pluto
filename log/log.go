package log

import (
	"io"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/pluto/config"
)

func GetLogger() log.Logger {
	var l log.Logger
	// get log config
	config := config.GetConfig().Log
	if config.Format.String() == "logfmt" {
		l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	} else {
		l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	}
	l = level.NewFilter(l, config.Level.GetLevelOption())
	l = log.With(l, "ts", log.DefaultTimestampUTC)
	if config.Level.String() == "debug" {
		l = log.With(l, "caller", log.DefaultCaller)
	}
	return l
}

func GetFileLogger(file io.Writer) log.Logger {
	var l log.Logger
	// get log config
	config := config.GetConfig().Log
	logFile := config.Path.String()
	if logFile != "" {

	}
	if config.Format.String() == "logfmt" {
		l = log.NewLogfmtLogger(log.NewSyncWriter(file))
	} else {
		l = log.NewJSONLogger(log.NewSyncWriter(file))
	}
	l = level.NewFilter(l, config.Level.GetLevelOption())
	l = log.With(l, "ts", log.DefaultTimestampUTC)
	if config.Level.String() == "debug" {
		l = log.With(l, "caller", log.DefaultCaller)
	}
	return l
}
