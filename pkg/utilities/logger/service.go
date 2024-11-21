package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type service struct {
	Log *logrus.Logger
}

var _ Service = (*service)(nil)
var once sync.Once
var l *logrus.Logger

func NewService(level string) *service {
	once.Do(func() {
		l = logrus.New()
		l.Level = loggerLevel(level)
	})
	return &service{
		Log: l,
	}
}

func (l service) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	l.Log.WithContext(ctx)
	l.Log.Info(msg, fields)
}

func (l service) Error(ctx context.Context, err error, fields map[string]interface{}) {
	l.Log.WithContext(ctx)
	l.Log.Error(err, fields)
}

func (l service) Debug(ctx context.Context, fields map[string]interface{}) {
	l.Log.WithContext(ctx)
	l.Log.Debug(fields)
}

func (l service) Warn(ctx context.Context, fields map[string]interface{}) {
	l.Log.WithContext(ctx)
	l.Log.Warn(fields)
}

func (l service) WrapError(err error, msg string) error {
	if err == nil {
		return errors.New(msg)
	}

	return fmt.Errorf("%w %s", err, msg)
}

func loggerLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
