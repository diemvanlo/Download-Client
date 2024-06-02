package database

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"goload/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	TabNameAccountPasswords           = "account_passwords"
	ColNameAccountPasswordOfAccountId = "of_user_id"
	ColNameAccountPasswordHash        = "hash"
)

type AccountPassword struct {
	OfUserID uint64 `sql:"of_user_id"`
	Hash     string `sql:"hash"`
}
type AccountPasswordDataAccessor interface {
	CreateUserPassword(ctx context.Context, accountPassword AccountPassword) error
	GetAccountPassword(ctx context.Context, ofAccountId uint64) (AccountPassword, error)
	UpdateUserPassword(ctx context.Context, accountPassword AccountPassword) error
	WithDatabase(database Database) AccountPasswordDataAccessor
}

type accountPasswordDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func (a accountPasswordDataAccessor) CreateUserPassword(ctx context.Context, accountPassword AccountPassword) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	_, err := a.database.
		Insert(TabNameAccountPasswords).
		Rows(goqu.Record{
			ColNameAccountPasswordOfAccountId: accountPassword.OfUserID,
			ColNameAccountPasswordHash:        accountPassword.Hash,
		}).
		Executor().ExecContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create account password")
		return status.Errorf(codes.Internal, "failed to create account password: %+v", err)
	}

	return nil
}

func (a accountPasswordDataAccessor) GetAccountPassword(ctx context.Context, ofAccountId uint64) (AccountPassword, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64(ColNameAccountPasswordOfAccountId, ofAccountId))
	accountPassword := AccountPassword{}
	found, err := a.database.From(TabNameAccountPasswords).Where(goqu.Ex{ColNameAccountPasswordOfAccountId: ofAccountId}).
		Where(goqu.Ex{ColNameAccountPasswordOfAccountId: ofAccountId}).ScanStructContext(ctx, &accountPassword)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get account password by id")
		return AccountPassword{}, status.Errorf(codes.Internal, "failed to get account password by id: %+v", err)
	}

	if !found {
		logger.Warn("Cannot find account by id")
		return AccountPassword{}, sql.ErrNoRows
	}

	return accountPassword, nil
}

func (a accountPasswordDataAccessor) UpdateUserPassword(ctx context.Context, accountPassword AccountPassword) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	_, err := a.database.Update(TabNameAccountPasswords).Set(goqu.Record{ColNameAccountPasswordHash: accountPassword.Hash}).
		Where(goqu.Ex{ColNameAccountPasswordOfAccountId: accountPassword.OfUserID}).Executor().ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to update account password")
		return status.Errorf(codes.Internal, "failed to update account password: %+v", err)
	}

	return nil
}

func (a accountPasswordDataAccessor) WithDatabase(database Database) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		database: database,
	}
}

func NewAccountPasswordDataAccessor(database *goqu.Database, logger *zap.Logger) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{database: database, logger: logger}
}
