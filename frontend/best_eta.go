package main

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/superliuwr/jaeger-demo/frontend/pkg/clients"
	"github.com/superliuwr/jaeger-demo/frontend/pkg/log"
	"github.com/superliuwr/jaeger-demo/frontend/pkg/pool"
)

const RouteWorkerPoolSize = 3

type bestETA struct {
	customer *clients.CustomerClient
	driver   *clients.DriverClient
	route    *clients.RouteClient
	pool     *pool.Pool
	logger   log.Factory
}

// Response contains ETA for a trip.
type Response struct {
	Driver string
	ETA    time.Duration
}

func newBestETA(tracer opentracing.Tracer, logger log.Factory, options ConfigOptions) *bestETA {
	return &bestETA{
		customer: clients.NewCustomerClient(
			tracer,
			logger.With(zap.String("component", "customer_client")),
			options.CustomerHostPort,
		),
		driver: clients.NewDriverClient(
			tracer,
			logger.With(zap.String("component", "driver_client")),
			options.DriverHostPort,
		),
		route: clients.NewRouteClient(
			tracer,
			logger.With(zap.String("component", "route_client")),
			options.RouteHostPort,
		),
		pool:   pool.New(RouteWorkerPoolSize),
		logger: logger,
	}
}

func (eta *bestETA) Get(ctx context.Context, customerID string) (*Response, error) {
	customer, err := eta.customer.GetCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	eta.logger.For(ctx).Info("Found customer", zap.Any("customer", customer))

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetBaggageItem("customer", customer.Name)
	}

	drivers, err := eta.driver.FindNearest(ctx, customer.Location)
	if err != nil {
		return nil, err
	}
	eta.logger.For(ctx).Info("Found drivers", zap.Any("drivers", drivers))

	results := eta.getRoutes(ctx, customer, drivers)
	eta.logger.For(ctx).Info("Found routes", zap.Any("routes", results))

	resp := &Response{ETA: math.MaxInt64}
	for _, result := range results {
		if result.err != nil {
			return nil, err
		}
		if result.route.ETA < resp.ETA {
			resp.ETA = result.route.ETA
			resp.Driver = result.driver
		}
	}
	if resp.Driver == "" {
		return nil, errors.New("no routes found")
	}

	eta.logger.For(ctx).Info("Dispatch successful", zap.String("driver", resp.Driver), zap.String("eta", resp.ETA.String()))
	return resp, nil
}

type routeResult struct {
	driver string
	route  *clients.Route
	err    error
}

// getRoutes calls Route service for each (customer, driver) pair
func (eta *bestETA) getRoutes(ctx context.Context, customer *clients.Customer, drivers []clients.Driver) []routeResult {
	results := make([]routeResult, 0, len(drivers))
	wg := sync.WaitGroup{}
	routesLock := sync.Mutex{}

	for _, dd := range drivers {
		wg.Add(1)
		driver := dd // capture loop var
		// Use worker pool to (potentially) execute requests in parallel
		eta.pool.Execute(func() {
			route, err := eta.route.FindRoute(ctx, driver.Location, customer.Location)
			routesLock.Lock()
			results = append(results, routeResult{
				driver: driver.DriverID,
				route:  route,
				err:    err,
			})
			routesLock.Unlock()
			wg.Done()
		})
	}

	wg.Wait()

	return results
}
