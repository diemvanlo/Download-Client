package http

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"goload/internal/configs"
	"goload/internal/generated/grpc/go_load"
	"goload/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	grpcConfig configs.GRPC
	httpConfig configs.HTTP
	logger     *zap.Logger
}

func NewServer(
	grpcConfig configs.GRPC,
	httpConfig configs.HTTP,
	logger *zap.Logger,
) Server {
	return &server{
		grpcConfig: grpcConfig,
		httpConfig: httpConfig,
		logger:     logger,
	}
}

func (s *server) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, s.logger)

	mux := runtime.NewServeMux()
	if err := go_load.RegisterGoLoadServiceHandlerFromEndpoint(
		ctx,
		mux,
		s.grpcConfig.Address,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}); err != nil {
		return err
	}

	httpServer := http.Server{
		Addr:              s.httpConfig.Address,
		ReadHeaderTimeout: time.Minute,
		Handler:           mux,
	}
	logger.With(zap.String("address", s.httpConfig.Address)).Info("starting http server")
	return httpServer.ListenAndServe()
}
