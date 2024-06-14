package grpc

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"go.uber.org/zap"
	"goload/internal/configs"
	"goload/internal/generated/go_load/v1"
	"goload/internal/utils"
	"google.golang.org/grpc"
	"net"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	handler    go_load.GoLoadServiceServer
	grpcConfig configs.GRPC
	logger     *zap.Logger
}

func NewServer(
	handler go_load.GoLoadServiceServer,
	grpcConfig configs.GRPC,
	logger *zap.Logger,
) Server {
	return &server{
		handler:    handler,
		grpcConfig: grpcConfig,
		logger:     logger,
	}
}

func (s *server) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, s.logger)

	listener, err := net.Listen("tcp", s.grpcConfig.Address)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to open tcp listener")
		return err
	}

	defer listener.Close()

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			validator.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
		),
	)
	go_load.RegisterGoLoadServiceServer(server, s.handler)

	logger.With(zap.String("address", s.grpcConfig.Address)).Info("starting grpc server")
	return server.Serve(listener)
}
