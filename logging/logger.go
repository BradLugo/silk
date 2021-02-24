package logging

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger interface {
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Error(name string, id string, err error)
	Warnf(name string, id string, msg string, args ...interface{})
	Infof(name string, id string, msg string, args ...interface{})
	Debugf(name string, id string, msg string, args ...interface{})
}

type WebberLogger struct {
	logger *zap.SugaredLogger
}

func NewWebberLogger() *WebberLogger {
	zl, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("error while instantiating the logging: %s", err.Error()))
	}
	defer func() {
		if err := zl.Sync(); err != nil {
			panic(fmt.Sprintf("could not sync logging %s", err.Error()))
		}
	}()

	sugar := zl.Sugar()

	return &WebberLogger{
		logger: sugar,
	}
}

func (wl *WebberLogger) Panicw(msg string, keysAndValues ...interface{}) {
	wl.logger.Panicw(msg, keysAndValues)
}

func (wl *WebberLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	wl.logger.Fatalw(msg, keysAndValues)
}

func (wl *WebberLogger) Errorw(msg string, keysAndValues ...interface{}) {
	wl.logger.Errorw(msg, keysAndValues)
}

func (wl *WebberLogger) Warnw(msg string, keysAndValues ...interface{}) {
	wl.logger.Warnw(msg, keysAndValues)
}

func (wl *WebberLogger) Infow(msg string, keysAndValues ...interface{}) {
	wl.logger.Infow(msg, keysAndValues)
}

func (wl *WebberLogger) Debugw(msg string, keysAndValues ...interface{}) {
	wl.logger.Debugw(msg, keysAndValues)
}

func (wl *WebberLogger) Error(name string, id string, err error) {
	wl.logger.Error(zap.String("name", name), zap.String("id", id), zap.Error(err))
}

func (wl *WebberLogger) Warnf(name string, id string, msg string, args ...interface{}) {
	wl.logger.Warnw(fmt.Sprintf(msg, args...), zap.String("name", name), zap.String("id", id))
}

func (wl *WebberLogger) Infof(name string, id string, msg string, args ...interface{}) {
	wl.logger.Infow(fmt.Sprintf(msg, args...), zap.String("name", name), zap.String("id", id))
}

func (wl *WebberLogger) Debugf(name string, id string, msg string, args ...interface{}) {
	wl.logger.Debugw(fmt.Sprintf(msg, args...), zap.String("name", name), zap.String("id", id))
}
