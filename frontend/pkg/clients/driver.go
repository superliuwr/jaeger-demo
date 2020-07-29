package clients

import (
	"context"
	"time"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/superliuwr/jaeger-demo/frontend/pkg/log"
)

// Driver describes a driver and the current car location.
type Driver struct {
	DriverID string
	Location string
}

type DriverClient struct {
	tracer opentracing.Tracer
	logger log.Factory
	client DriverServiceClient
}

// NewDriverClient creates a new driver.Client
func NewDriverClient(tracer opentracing.Tracer, logger log.Factory, hostPort string) *DriverClient {
	conn, err := grpc.Dial(hostPort, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(tracer)))
	if err != nil {
		logger.Bg().Fatal("Cannot create gRPC connection", zap.Error(err))
	}

	client := NewDriverServiceClient(conn)

	return &DriverClient{
		tracer: tracer,
		logger: logger,
		client: client,
	}
}

// FindNearest implements driver.Interface#FindNearest as an RPC
func (c *DriverClient) FindNearest(ctx context.Context, location string) ([]Driver, error) {
	c.logger.For(ctx).Info("Finding nearest drivers", zap.String("location", location))
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	response, err := c.client.FindNearest(ctx, &DriverLocationRequest{Location: location})
	if err != nil {
		return nil, err
	}

	return fromProto(response), nil
}

func fromProto(response *DriverLocationResponse) []Driver {
	retMe := make([]Driver, len(response.Locations))
	for i, result := range response.Locations {
		retMe[i] = Driver{
			DriverID: result.DriverID,
			Location: result.Location,
		}
	}

	return retMe
}