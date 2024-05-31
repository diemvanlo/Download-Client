package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"goload/internal/dataAccess/database"
)

type CreateAccountParams struct {
	Username string
	Password string
}

type CreateAccountOutput struct {
	ID       uint64
	UserName string
}

type CreateSessionParams struct {
	Username string
	Password string
}

type Account interface {
	CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (string, error)
}

type account struct {
	goquDatabase                *goqu.Database
	accountDataAccessor         database.AccountDataAccessor
	accountPasswordDataAccessor database.AccountPasswordDataAccessor
	hashLogic                   Hash
	//tokenLogic Token
}

func NewAccount(
	goquDatabase *goqu.Database,
	accountDataAccessor database.AccountDataAccessor,
	accountPasswordDataAccessor database.AccountPasswordDataAccessor,
	hashLogic Hash,
) Account {
	return &account{
		goquDatabase:                goquDatabase,
		accountDataAccessor:         accountDataAccessor,
		accountPasswordDataAccessor: accountPasswordDataAccessor,
		hashLogic:                   hashLogic,
	}
}

func (a account) isAccontUsernameTaken(ctx context.Context, username string) (bool, error) {
	_, err := a.accountDataAccessor.GetAccountByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a account) CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error) {
	var accountId uint64

	txErr := a.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		usernameTaken, err := a.isAccontUsernameTaken(ctx, params.Username)
		if err != nil {
			return err
		}

		if usernameTaken {
			return errors.New("Username is already taken")
		}
		accountId, err = a.accountDataAccessor.WithDatabase(td).CreateAccount(ctx, database.Account{Username: params.Username})
		if err != nil {
			return err
		}

		hashedPassword, err := a.hashLogic.Hash(ctx, params.Password)
		if err != nil {
			return err
		}

		if err := a.accountPasswordDataAccessor.WithDatabase(td).CreateUserPassword(ctx, database.AccountPassword{
			OfUserID: accountId,
			Hash:     hashedPassword,
		}); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return CreateAccountOutput{}, txErr
	}
	return CreateAccountOutput{ID: accountId, UserName: params.Username}, nil
}

func (a account) CreateSession(ctx context.Context, params CreateSessionParams) (string, error) {
	existingAccount, err := a.accountDataAccessor.GetAccountByUsername(ctx, params.Username)
	if err != nil {
		return "", err
	}
	existingAccountPassword, err := a.accountPasswordDataAccessor.GetAccountPassword(ctx, existingAccount.ID)
	if err != nil {
		return "", err
	}
	isHashEqual, err := a.hashLogic.IsHashEqual(ctx, params.Password, existingAccountPassword.Hash)
	if err != nil {
		return "", err
	}
	if isHashEqual {
		return "", errors.New("incorrect password")
	}
	return "", nil
}
