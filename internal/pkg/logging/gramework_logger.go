package logging

import (
	"fmt"

	apexlog "github.com/apex/log"
	"github.com/gramework/gramework"
	"go.uber.org/zap"
)

func NewGrameworkLogger(log *zap.Logger) apexlog.Interface {
	instance := &grameworkLogger{
		log: log,
	}
	instance.apexLog = &apexlog.Logger{
		Level:   apexlog.DebugLevel,
		Handler: instance,
	}

	gramework.Logger.Handler = instance

	return instance
}

type grameworkLogger struct {
	log     *zap.Logger
	apexLog *apexlog.Logger
}

func (impl *grameworkLogger) HandleLog(entry *apexlog.Entry) (err error) {
	var logFunc func(msg string, fields ...zap.Field)
	switch entry.Level {
	case apexlog.FatalLevel:
		logFunc = impl.log.Fatal
	case apexlog.ErrorLevel:
		logFunc = impl.log.Error
	case apexlog.WarnLevel:
		logFunc = impl.log.Warn
	case apexlog.InfoLevel:
		logFunc = impl.log.Info
	case apexlog.DebugLevel:
		logFunc = impl.log.Debug
	default:
		impl.log.Warn("can't define log level")
		return
	}

	keys := entry.Fields.Names()
	zapFields := make([]zap.Field, 0, len(keys))
	for _, key := range keys {
		zapFields = append(zapFields, zap.Any(key, entry.Fields.Get(key)))
	}

	logFunc(entry.Message, zapFields...)

	return
}

func (impl *grameworkLogger) WithFields(fields apexlog.Fielder) *apexlog.Entry {
	return apexlog.NewEntry(impl.apexLog).WithFields(fields)
}

func (impl *grameworkLogger) WithField(key string, value interface{}) *apexlog.Entry {
	return apexlog.NewEntry(impl.apexLog).WithField(key, value)
}

func (impl *grameworkLogger) WithError(err error) *apexlog.Entry {
	return apexlog.NewEntry(impl.apexLog).WithError(err)
}

func (impl *grameworkLogger) Debug(msg string) {
	impl.log.Debug(msg)
}

func (impl *grameworkLogger) Info(msg string) {
	impl.log.Info(msg)
}

func (impl *grameworkLogger) Warn(msg string) {
	impl.log.Warn(msg)
}

func (impl *grameworkLogger) Error(msg string) {
	impl.log.Error(msg)
}

func (impl *grameworkLogger) Fatal(msg string) {
	impl.log.Fatal(msg)
}

func (impl *grameworkLogger) Debugf(msg string, v ...interface{}) {
	impl.log.Debug(fmt.Sprintf(msg, v...))
}

func (impl *grameworkLogger) Infof(msg string, v ...interface{}) {
	impl.log.Info(fmt.Sprintf(msg, v...))
}

func (impl *grameworkLogger) Warnf(msg string, v ...interface{}) {
	impl.log.Warn(fmt.Sprintf(msg, v...))
}

func (impl *grameworkLogger) Errorf(msg string, v ...interface{}) {
	impl.log.Error(fmt.Sprintf(msg, v...))
}

func (impl *grameworkLogger) Fatalf(msg string, v ...interface{}) {
	impl.log.Fatal(fmt.Sprintf(msg, v...))
}

func (impl *grameworkLogger) Trace(msg string) *apexlog.Entry {
	return apexlog.NewEntry(impl.apexLog).Trace(msg)
}
