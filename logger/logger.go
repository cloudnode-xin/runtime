package logger

import (
	"context"
	"strings"

	"github.com/cloudnodexin/runtime"
	"github.com/sirupsen/logrus"
)

type setup struct {
	options []runtime.LoggerOption
}

func (s *setup) Name() string {
	return "#loggersetup"
}

func (s *setup) IsHealthy() bool {
	return true
}

func (s *setup) Load(f runtime.Finder) error {
	log := f.MustGet("#logger").(*runtime.Logger)

	for _, opt := range s.options {
		log.Setup(opt)
	}

	return nil
}

func (s *setup) Start(f runtime.Finder, ctx context.Context) error {
	return nil
}

func (s *setup) Stop(f runtime.Finder) error {
	return nil
}

func Setup(opts ...runtime.LoggerOption) runtime.Servicer {
	return &setup{
		options: opts,
	}
}

func FormatString(format string) runtime.LoggerOption {
	var formatter logrus.Formatter

	formatter = &logrus.TextFormatter{}
	if strings.ToLower(format) == "json" {
		formatter = &logrus.JSONFormatter{}
	}

	return Format(formatter)
}

func Format(formatter logrus.Formatter) runtime.LoggerOption {
	return func(log *logrus.Logger) {
		log.SetFormatter(formatter)
	}
}

func LevelString(level string) runtime.LoggerOption {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.WarnLevel
	}

	return Level(lvl)
}

func Level(level logrus.Level) runtime.LoggerOption {
	return func(log *logrus.Logger) {
		log.SetLevel(level)
	}
}
