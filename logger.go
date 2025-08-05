package runtime

import (
	"context"

	"github.com/sirupsen/logrus"
)

type LoggerOption func(log *logrus.Logger)

type Logger struct {
	log *logrus.Logger
}

func (l *Logger) Name() string {
	return "#logger"
}

func (l *Logger) IsHealthy() bool {
	return true
}

func (l *Logger) Load(f Finder) error {
	return nil
}

func (l *Logger) Start(f Finder, ctx context.Context) error {
	return nil
}

func (l *Logger) Stop(f Finder) error {
	return nil
}

func (l *Logger) New(target string) *logrus.Entry {
	return l.log.WithField("target", target)
}

func (l *Logger) Setup(opt LoggerOption) {
	opt(l.log)
}

func logger() Servicer {
	l := &Logger{}
	l.log = logrus.New()
	l.log.SetLevel(logrus.TraceLevel)

	return l
}
