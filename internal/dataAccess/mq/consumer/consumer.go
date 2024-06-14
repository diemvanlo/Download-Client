package consumer

import (
	"context"
	"fmt"
	"goload/internal/configs"
	"goload/internal/utils"

	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type HandlerFunc func(ctx context.Context, queueName string, payload []byte) error

type Consumer interface {
	RegisterHandler(queueName string, handlerFunc HandlerFunc)
	Start(ctx context.Context) error
}

type consumer struct {
	saramaConsumer            sarama.Consumer
	queueNameToHandlerFuncMap map[string]HandlerFunc
	logger                    *zap.Logger
}

func newSaramaConfig(mqConfig configs.MQ) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = mqConfig.ClientID
	saramaConfig.Metadata.Full = true
	return saramaConfig
}

func NewConsumer(
	mqConfig configs.MQ,
	logger *zap.Logger,
) (Consumer, error) {
	saramaConsumer, err := sarama.NewConsumer(mqConfig.Addresses, newSaramaConfig(mqConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama consumer: %w", err)
	}

	return &consumer{
		saramaConsumer:            saramaConsumer,
		logger:                    logger,
		queueNameToHandlerFuncMap: make(map[string]HandlerFunc),
	}, nil
}

func (c *consumer) RegisterHandler(queueName string, handlerFunc HandlerFunc) {
	c.queueNameToHandlerFuncMap[queueName] = handlerFunc
}

func (c consumer) consumer(queueName string, handlerFunc HandlerFunc, exitSignalChannel chan os.Signal) error {
	logger := c.logger.With(zap.String("queue_name", queueName))

	partitionConsumer, err := c.saramaConsumer.ConsumePartition(queueName, 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("failed to create sarama partition consumer: %w", err)
	}

	for {
		select {
		case message := <-partitionConsumer.Messages():
			if err := handlerFunc(context.Background(), queueName, message.Value); err != nil {
				logger.With(zap.Error(err)).Error("failed to handle message")
			}
		case <-exitSignalChannel:
			return nil
		}
	}
}

func (c consumer) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, c.logger)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for queueName, handlerFunc := range c.queueNameToHandlerFuncMap {
		go func(queueName string, handlerFunc HandlerFunc) {
			if err := c.consumer(queueName, handlerFunc, signals); err != nil {
				logger.With(zap.String("queue_name", queueName)).
					With(zap.Error(err)).
					Error("failed to consume")
			}
		}(queueName, handlerFunc)
	}

	<-signals
	return nil
}
