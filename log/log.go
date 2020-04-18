package log

import (
	"io"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/pluto/config"
)

type PlutoLog struct {
	logger log.Logger
}

func (pl *PlutoLog) Error(message interface{}) {
	level.Error(pl.logger).Log("error", message)
}

func (pl *PlutoLog) Info(message interface{}) {
	level.Info(pl.logger).Log("info", message)
}

func (pl *PlutoLog) Debug(message interface{}) {
	level.Debug(pl.logger).Log("debug", message)
}

func (pl *PlutoLog) Warn(message interface{}) {
	level.Debug(pl.logger).Log("warning", message)
}

func (pl *PlutoLog) With(keyvals ...interface{}) *PlutoLog {
	l := log.With(pl.logger, keyvals...)
	return &PlutoLog{
		logger: l,
	}
}

func NewLogger(config *config.Config) (*PlutoLog, error) {
	var l log.Logger
	var writer io.WriteCloser
	var err error
	if *config.Log.File == "" {
		writer = os.Stdout
	} else {
		writer, err = os.OpenFile(*config.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}
	if config.Log.Format.String() == "logfmt" {
		l = log.NewLogfmtLogger(log.NewSyncWriter(writer))
	} else {
		l = log.NewJSONLogger(log.NewSyncWriter(writer))
	}
	l = level.NewFilter(l, config.Log.Level.GetLevelOption())
	l = log.With(l, "ts", log.DefaultTimestamp)
	return &PlutoLog{
		logger: l,
	}, nil
}
