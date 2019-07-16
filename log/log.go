package log

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/pluto/config"
)

func NewLogger(config *config.LogConfig) log.Logger {
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
