package logging

import (
	stdlog "log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logInstance *zap.Logger

func init() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "M",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     false,
		DisableStacktrace: false,
	}

	log, err := cfg.Build(zap.AddStacktrace(zap.WarnLevel))
	if err != nil {
		stdlog.Fatal(err)
	}
	logInstance = log
}

func New(name string) *zap.Logger {
	if len(name) == 0 {
		return logInstance
	}

	return logInstance.Named(name)
}
