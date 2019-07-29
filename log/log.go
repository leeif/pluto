package log

import (
	"io"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/pluto/config"
)

func GetLogger(config *config.LogConfig) log.Logger {
	var l log.Logger
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

func GetFileLogger(config *config.LogConfig, file io.Writer) log.Logger {
	var l log.Logger
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
