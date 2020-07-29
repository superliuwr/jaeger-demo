package main

import (
	"net"
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/superliuwr/jaeger-demo/driver/pkg/log"
	"github.com/superliuwr/jaeger-demo/driver/pkg/tracing"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(-1)
	}
}

func execute() error {
	rootLogger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)
	appLogger := rootLogger.With(zap.String("service", "driver"))
	loggerFactory := log.NewFactory(appLogger)

	server := NewServer(
		net.JoinHostPort("0.0.0.0", strconv.Itoa(8081)),
		tracing.Init("driver", loggerFactory),
		loggerFactory,
	)

	return logError(appLogger, server.Run())
}

func logError(logger *zap.Logger, err error) error {
	if err != nil {
		logger.Error("Error running command", zap.Error(err))
	}

	return err
}