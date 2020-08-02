package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	"github.com/superliuwr/jaeger-demo/driver/delay"
	"github.com/superliuwr/jaeger-demo/driver/log"
	"github.com/superliuwr/jaeger-demo/driver/tracing"
)

var (
	// RedisFindDelay is how long finding closest drivers takes.
	RedisFindDelay = 20 * time.Millisecond

	// RedisFindDelayStdDev is standard deviation.
	RedisFindDelayStdDev = RedisFindDelay / 4

	// RedisGetDelay is how long retrieving a driver record takes.
	RedisGetDelay = 10 * time.Millisecond

	// RedisGetDelayStdDev is standard deviation
	RedisGetDelayStdDev = RedisGetDelay / 4
)

// Redis is a simulator of remote Redis cache
type Redis struct {
	tracer opentracing.Tracer // simulate redis as a separate process
	logger log.Factory
	errorSimulator
}

func newRedis(logger log.Factory) *Redis {
	return &Redis{
		tracer: tracing.Init("redis", logger),
		logger: logger,
	}
}

// FindDriverIDs finds IDs of drivers who are near the location.
func (r *Redis) FindDriverIDs(ctx context.Context, location string) []string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := r.tracer.StartSpan("FindDriverIDs", opentracing.ChildOf(span.Context()))

		span.SetTag("param.location", location)
		ext.SpanKindRPCClient.Set(span)
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	// simulate RPC delay
	delay.Sleep(RedisFindDelay, RedisFindDelayStdDev)

	drivers := make([]string, 10)
	for i := range drivers {
		// #nosec
		drivers[i] = fmt.Sprintf("T7%05dC", rand.Int()%100000)
	}
	r.logger.For(ctx).Info("Found drivers", zap.Strings("drivers", drivers))

	return drivers
}

// GetDriver returns driver and the current car location
func (r *Redis) GetDriver(ctx context.Context, driverID string) (Driver, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := r.tracer.StartSpan("GetDriver", opentracing.ChildOf(span.Context()))

		span.SetTag("param.driverID", driverID)
		ext.SpanKindRPCClient.Set(span)
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	// simulate RPC delay
	delay.Sleep(RedisGetDelay, RedisGetDelayStdDev)

	if err := r.checkError(); err != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			ext.Error.Set(span, true)
		}

		r.logger.For(ctx).Error("redis timeout", zap.String("driver_id", driverID), zap.Error(err))

		return Driver{}, err
	}

	// #nosec
	return Driver{
		DriverID: driverID,
		Location: fmt.Sprintf("%d,%d", rand.Int()%1000, rand.Int()%1000),
	}, nil
}

var errTimeout = errors.New("redis timeout")

type errorSimulator struct {
	sync.Mutex
	countTillError int
}

func (es *errorSimulator) checkError() error {
	es.Lock()
	es.countTillError--

	if es.countTillError > 0 {
		es.Unlock()
		return nil
	}

	es.countTillError = 5
	es.Unlock()

	delay.Sleep(2*RedisGetDelay, 0) // add more delay for "timeout"

	return errTimeout
}
