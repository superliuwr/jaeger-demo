package tracing

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"

	"github.com/superliuwr/jaeger-demo/frontend/log"
)

// Init creates a new instance of Jaeger tracer.
func Init(serviceName string, logger log.Factory) opentracing.Tracer {
	// Read host and port from Env Vars
	cfg, err := config.FromEnv()
	if err != nil {
		logger.Bg().Fatal("cannot parse Jaeger env vars", zap.Error(err))
	}

	cfg.ServiceName = serviceName
	// Always sample all requests
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1

	jaegerLogger := jaegerLoggerAdapter{logger.Bg()}

	tracer, _, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
	)
	if err != nil {
		logger.Bg().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	}

	return tracer
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
