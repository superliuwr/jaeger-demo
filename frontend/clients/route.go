package clients

import (
	"context"
	"net/http"
	"net/url"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/superliuwr/jaeger-demo/frontend/log"
	"github.com/superliuwr/jaeger-demo/frontend/tracing"
)

// Route describes a route between Pickup and Dropoff locations and expected time to arrival.
type Route struct {
	Pickup  string
	Dropoff string
	ETA     int
}

type RouteClient struct {
	tracer   opentracing.Tracer
	logger   log.Factory
	client   *tracing.HTTPClient
	hostPort string
}

// NewRouteClient creates a new route.Client
func NewRouteClient(tracer opentracing.Tracer, logger log.Factory, hostPort string) *RouteClient {
	return &RouteClient{
		tracer: tracer,
		logger: logger,
		client: &tracing.HTTPClient{
			Client: &http.Client{Transport: &nethttp.Transport{}},
			Tracer: tracer,
		},
		hostPort: hostPort,
	}
}

// FindRoute implements route.Interface#FindRoute as an RPC
func (c *RouteClient) FindRoute(ctx context.Context, pickup, dropoff string) (*Route, error) {
	c.logger.For(ctx).Info("Finding route", zap.String("pickup", pickup), zap.String("dropoff", dropoff))

	v := url.Values{}
	v.Set("pickup", pickup)
	v.Set("dropoff", dropoff)
	url := "http://" + c.hostPort + "/route?" + v.Encode()

	var route Route

	if err := c.client.GetJSON(ctx, "/route", url, &route); err != nil {
		c.logger.For(ctx).Error("Error getting route", zap.Error(err))

		return nil, err
	}

	return &route, nil
}
