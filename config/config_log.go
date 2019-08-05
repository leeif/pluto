package config

import (
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

type LogConfig struct {
	Level  *AllowedLevel  `kiper_value:"name:level;help:log level = debug, info, warn, error;default:info"`
	Format *AllowedFormat `kiper_value:"name:format;help:log format = json, logfmt;default:logfmt"`
	File   *string        `kiper_value:"name:file;help:log file path;default:./pluto.log"`
}

type AllowedLevel struct {
	s string
	o level.Option
}

func (l *AllowedLevel) String() string {
	return l.s
}

func (l *AllowedLevel) Set(s string) error {
	switch s {
	case "debug":
		l.o = level.AllowDebug()
	case "info":
		l.o = level.AllowInfo()
	case "warn":
		l.o = level.AllowWarn()
	case "error":
		l.o = level.AllowError()
	default:
		return errors.Errorf("unrecognized log level %q", s)
	}
	l.s = s
	return nil
}

func (l *AllowedLevel) GetLevelOption() level.Option {
	return l.o
}

type AllowedFormat struct {
	s string
}

func (f *AllowedFormat) Set(s string) error {
	switch s {
	case "logfmt", "json":
		f.s = s
	default:
		return errors.Errorf("unrecognized log format %q", s)
	}
	return nil
}

func (f *AllowedFormat) String() string {
	return f.s
}

func newLogConfig() *LogConfig {
	return &LogConfig{
		Level:  &AllowedLevel{},
		Format: &AllowedFormat{},
	}
}
