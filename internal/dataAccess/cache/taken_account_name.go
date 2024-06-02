package cache

import (
	"context"
	"goload/internal/utils"

	"go.uber.org/zap"
)

const (
	setKeyNameTakenAccountName = "taken_account_name_set"
)

type TakenAccountName interface {
	Add(ctx context.Context, accountName string) error
	Has(ctx context.Context, accountName string) (bool, error)
}

type takenAccountName struct {
	client Client
	logger *zap.Logger
}

func NewTakenAccountName(
	client Client,
	logger *zap.Logger,
) TakenAccountName {
	return &takenAccountName{
		client: client,
		logger: logger,
	}
}

func (c takenAccountName) Add(ctx context.Context, accountName string) error {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("account_name", accountName))
	if err := c.client.AddToSet(ctx, setKeyNameTakenAccountName, accountName); err != nil {
		logger.With(zap.Error(err)).Error("failed to add username to set in cache")
		return err
	}

	return nil
}

func (c takenAccountName) Has(ctx context.Context, accountName string) (bool, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("account_name", accountName))
	result, err := c.client.IsDataSet(ctx, setKeyNameTakenAccountName, accountName)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to check if account name is in set in cache")
	}

	return result, nil
}
