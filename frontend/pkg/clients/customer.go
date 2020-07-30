package clients

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/superliuwr/jaeger-demo/frontend/pkg/log"
	"github.com/superliuwr/jaeger-demo/frontend/pkg/tracing"
)

// Customer contains data about a customer.
type Customer struct {
	ID       string
	Name     string
	Location string
}

type CustomerClient struct {
	tracer   opentracing.Tracer
	logger   log.Factory
	client   *tracing.HTTPClient
	hostPort string
}

// NewCustomerClient creates a new customer.Client
func NewCustomerClient(tracer opentracing.Tracer, logger log.Factory, hostPort string) *CustomerClient {
	return &CustomerClient{
		tracer: tracer,
		logger: logger,
		client: &tracing.HTTPClient{
			Client: &http.Client{Transport: &nethttp.Transport{}},
			Tracer: tracer,
		},
		hostPort: hostPort,
	}
}

// GetCustomer implements customer.Interface#Get as an RPC
func (c *CustomerClient) GetCustomer(ctx context.Context, customerID string) (*Customer, error) {
	c.logger.For(ctx).Info("Getting customer", zap.String("customer_id", customerID))

	url := fmt.Sprintf("http://"+c.hostPort+"/customer?customer=%s", customerID)
	fmt.Println(url)

	var customer Customer
	if err := c.client.GetJSON(ctx, "/customer", url, &customer); err != nil {
		return nil, err
	}

	return &customer, nil
}
