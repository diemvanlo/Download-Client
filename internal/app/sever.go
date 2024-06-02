package app

import (
	"context"
	"go.uber.org/zap"
	"goload/internal/handler/grpc"
	"goload/internal/handler/http"
	"goload/internal/utils"
	"syscall"
)

type Server struct {
	grpcSever grpc.Server
	httpSever http.Server
	logger    *zap.Logger
}

func NewSever(grpcServer grpc.Server, httpServer http.Server, logger *zap.Logger) *Server {
	return &Server{logger: logger, grpcSever: grpcServer, httpSever: httpServer}
}

func (s Server) Start() error {
	go func() {
		err := s.grpcSever.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("grpc sever stopped")
	}()

	go func() {
		err := s.httpSever.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("http sever stopped")
	}()

	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
	return nil
}
