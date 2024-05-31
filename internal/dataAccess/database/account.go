package database

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"goload/internal/utils"
)

const (
	TabNameAccounts    = "accounts"
	ColNameAccountsID  = "id"
	ColNameAccountName = "account_name"
)

type Account struct {
	ID       uint64 `sql:"id"`
	Username string `sql:"account_name"`
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account Account) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByUsername(ctx context.Context, username string) (Account, error)
	WithDatabase(database Database) AccountDataAccessor
}

type accountDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func (a accountDataAccessor) CreateAccount(ctx context.Context, account Account) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	// need to get back
	result, err := a.database.Insert(TabNameAccounts).Rows(goqu.Record{
		ColNameAccountName: account.Username,
	}).Executor().ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("Fail to create account")
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get last inserted id")
		return 0, err
	}

	return uint64(lastInsertID), nil
}

func NewDatabaseAccessor(database *goqu.Database, logger *zap.Logger) AccountDataAccessor {
	return &accountDataAccessor{database: database, logger: logger}
}

func (a accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)
	account := Account{}
	found, err := a.database.From(TabNameAccounts).Where(goqu.Ex{ColNameAccountsID: id}).ScanStructContext(ctx, &account)

	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get account by id")
		return Account{}, err
	}

	if !found {
		logger.Warn("cannot find account by id")
	}

	return account, nil
}

func (a accountDataAccessor) GetAccountByUsername(ctx context.Context, username string) (Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)
	account := Account{}
	found, err := a.database.From(TabNameAccounts).Where(goqu.Ex{ColNameAccountName: username}).ScanStructContext(ctx, &account)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get account by id")
		return Account{}, err
	}

	if !found {
		logger.Warn("cannot find account by id")
		return Account{}, err
	}

	return account, nil
}

func (a accountDataAccessor) WithDatabase(database Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
	}
}
