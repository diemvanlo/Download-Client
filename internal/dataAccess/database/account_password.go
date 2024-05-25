package database

import (
	"context"
	"github.com/doug-martin/goqu/v9"
)

type AccountPassword struct {
	OfUserID uint64 `sql:"of_user_id"`
	Hash     string `sql:"hash"`
}
type AccountPasswordDataAccessor interface {
	CreateUserPassword(ctx context.Context, accountPassword AccountPassword) error
	UpdateUserPassword(ctx context.Context, accountPassword AccountPassword) error
	WithDatabase(database Database) AccountPasswordDataAccessor
}

type accountPasswordDataAccessor struct {
	database Database
}

func NewAccountPasswordDataAccessor(database *goqu.Database) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{database: database}
}

func (a accountPasswordDataAccessor) CreateUserPassword(ctx context.Context, accountPassword AccountPassword) error {
	//TODO implement me
	panic("implement me")
}

func (a accountPasswordDataAccessor) UpdateUserPassword(ctx context.Context, accountPassword AccountPassword) error {
	//TODO implement me
	panic("implement me")
}

func (a accountPasswordDataAccessor) WithDatabase(database Database) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		database: database,
	}
}
