package main

import (
	"context"
	"net"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/superliuwr/jaeger-demo/driver/log"
)

// Driver describes a driver and the current car location.
type Driver struct {
	DriverID string
	Location string
}

type Server struct {
	hostPort string
	tracer   opentracing.Tracer
	logger   log.Factory
	redis    *Redis
	server   *grpc.Server
}

var _ DriverServiceServer = (*Server)(nil)

// NewServer creates a new driver.Server
func NewServer(hostPort string, tracer opentracing.Tracer, logger log.Factory) *Server {
	server := grpc.NewServer(grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(tracer)),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracer)))

	return &Server{
		hostPort: hostPort,
		tracer:   tracer,
		logger:   logger,
		server: server,
		redis:    newRedis(logger),
	}
}

// Run starts the Driver server
func (s *Server) Run() error {
	s.logger.Bg().Info("Starting", zap.String("address", "http://"+s.hostPort))

	lis, err := net.Listen("tcp", s.hostPort)
	if err != nil {
		s.logger.Bg().Fatal("Unable to create http listener", zap.Error(err))
	}

	RegisterDriverServiceServer(s.server, s)

	err = s.server.Serve(lis)
	if err != nil {
		s.logger.Bg().Fatal("Unable to start gRPC server", zap.Error(err))
	}

	return err
}

// FindNearest implements gRPC driver interface
func (s *Server) FindNearest(ctx context.Context, location *DriverLocationRequest) (*DriverLocationResponse, error) {
	s.logger.For(ctx).Info("Searching for nearby drivers", zap.String("location", location.Location))
	driverIDs := s.redis.FindDriverIDs(ctx, location.Location)

	retMe := make([]*DriverLocation, len(driverIDs))
	for i, driverID := range driverIDs {
		var drv Driver
		var err error

		for i := 0; i < 3; i++ {
			drv, err = s.redis.GetDriver(ctx, driverID)
			if err == nil {
				break
			}
			s.logger.For(ctx).Error("Retrying GetDriver after error", zap.Int("retry_no", i+1), zap.Error(err))
		}
		if err != nil {
			s.logger.For(ctx).Error("Failed to get driver after 3 attempts", zap.Error(err))
			return nil, err
		}

		retMe[i] = &DriverLocation{
			DriverID: drv.DriverID,
			Location: drv.Location,
		}
	}

	s.logger.For(ctx).Info("Search successful", zap.Int("num_drivers", len(retMe)))

	return &DriverLocationResponse{Locations: retMe}, nil
}
