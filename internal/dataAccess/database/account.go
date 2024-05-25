package database

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"log"
)

type Account struct {
	UserID   uint64 `sql:"user_id"`
	Username string `sql:"username"`
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account Account) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByUsername(ctx context.Context, username string) (Account, error)
	WithDatabase(database Database) AccountDataAccessor
}

type accountDataAccessor struct {
	database Database
}

func (a accountDataAccessor) CreateAccount(ctx context.Context, account Account) (uint64, error) {
	result, err := a.database.Insert("accounts").Rows(goqu.Record{
		"username": account.Username,
	}).Executor().ExecContext(ctx)

	if err != nil {
		log.Printf("Failed to create account, err=%+v", err)
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last inserted id, error=%+v\n", err)
		return 0, err
	}

	return uint64(lastInsertID), nil
}

func NewDatabaseAccessor(database *goqu.Database) AccountDataAccessor {
	return &accountDataAccessor{database: database}
}

func (a accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a accountDataAccessor) GetAccountByUsername(ctx context.Context, username string) (Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a accountDataAccessor) WithDatabase(database Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
	}
}
