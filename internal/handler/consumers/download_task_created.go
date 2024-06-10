package consumers

import (
	"context"
	"goload/internal/dataAccess/mq/producer"
	"goload/internal/utils"

	"go.uber.org/zap"
)

type DownloadTaskCreated interface {
	Handle(ctx context.Context, event producer.DownloadTaskCreated) error
}

type downloadTaskCreated struct {
	logger *zap.Logger
}

func NewDownloadTaskCreated(
	logger *zap.Logger,
) DownloadTaskCreated {
	return &downloadTaskCreated{
		logger: logger,
	}
}

func (d downloadTaskCreated) Handle(ctx context.Context, event producer.DownloadTaskCreated) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("event", event))
	logger.Info("download task created event received")

	return nil
}