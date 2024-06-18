package app

import (
	"context"
	"goload/internal/handler/consumers"
	"goload/internal/handler/grpc"
	"goload/internal/handler/http"
	"goload/internal/utils"
	"syscall"

	"go.uber.org/zap"
)

type Server struct {
	grpcServer   grpc.Server
	httpServer   http.Server
	rootConsumer consumers.Root
	logger       *zap.Logger
}

func NewServer(
	grpcServer grpc.Server,
	httpServer http.Server,
	rootConsumer consumers.Root,
	logger *zap.Logger,
) *Server {
	return &Server{
		grpcServer:   grpcServer,
		httpServer:   httpServer,
		rootConsumer: rootConsumer,
		logger:       logger,
	}
}

func (s Server) Start() error {
	go func() {
		err := s.grpcServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("grpc server stopped")
	}()

	go func() {
		err := s.httpServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("http server stopped")
	}()

	go func() {
		err := s.rootConsumer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("message queue consumer stopped")
	}()

	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
	return nil
}
