package main

import (
	"net"
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/superliuwr/jaeger-demo/frontend/pkg/log"
	"github.com/superliuwr/jaeger-demo/frontend/pkg/tracing"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(-1)
	}
}

func execute() error {
	var options ConfigOptions

	options.FrontendHostPort = net.JoinHostPort("0.0.0.0", strconv.Itoa(8080))
	options.DriverHostPort = net.JoinHostPort("driver", strconv.Itoa(8081))
	options.CustomerHostPort = net.JoinHostPort("0.0.0.0", strconv.Itoa(8082))
	options.RouteHostPort = net.JoinHostPort("0.0.0.0", strconv.Itoa(8083))
	options.BasePath = `/`

	rootLogger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)
	appLogger := rootLogger.With(zap.String("service", "frontend"))
	loggerFactory := log.NewFactory(appLogger)

	server := NewServer(
		options,
		tracing.Init("frontend", loggerFactory),
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