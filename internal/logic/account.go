package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"goload/internal/dataAccess/cache"
	"goload/internal/dataAccess/database"
	"goload/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	takenAccountNameCache       cache.TakenAccountName
	accountDataAccessor         database.AccountDataAccessor
	accountPasswordDataAccessor database.AccountPasswordDataAccessor
	hashLogic                   Hash
	tokenLogic                  Token
	logger                      *zap.Logger
}

func NewAccount(
	goquDatabase *goqu.Database,
	takenAccountName cache.TakenAccountName,
	accountDataAccessor database.AccountDataAccessor,
	accountPasswordDataAccessor database.AccountPasswordDataAccessor,
	hashLogic Hash,
	tokenLogic Token,
	logger *zap.Logger,
) Account {
	return &account{
		goquDatabase:                goquDatabase,
		takenAccountNameCache:       takenAccountName,
		accountDataAccessor:         accountDataAccessor,
		accountPasswordDataAccessor: accountPasswordDataAccessor,
		hashLogic:                   hashLogic,
		tokenLogic:                  tokenLogic,
		logger:                      logger,
	}
}

func (a account) isAccontUsernameTaken(ctx context.Context, username string) (bool, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("username", username))

	accountNameTaken, err := a.takenAccountNameCache.Has(ctx, username)
	if err != nil {
		logger.With(zap.Error(err)).Warn("Failed to get username in cache, will fall back to database")
	} else {
		return accountNameTaken, nil
	}

	_, err = a.accountDataAccessor.GetAccountByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if err := a.takenAccountNameCache.Add(ctx, username); err != nil {
		logger.With(zap.Error(err)).Warn("failed to set username into taken set in cache")
	}
	return true, nil
}

func (a account) CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error) {
	accountNameTaken, err := a.isAccontUsernameTaken(ctx, params.Username)
	if err != nil {
		return CreateAccountOutput{}, status.Errorf(codes.Internal, "failed to check if username is taken")
	}

	if accountNameTaken {
		return CreateAccountOutput{}, status.Errorf(codes.AlreadyExists, "account name is already taken")
	}

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

		hashedPassword, hashErr := a.hashLogic.Hash(ctx, params.Password)
		if hashErr != nil {
			return hashErr
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
		return "", status.Errorf(codes.Unauthenticated, "incorrect password")
	}

	token, _, err := a.tokenLogic.GetToken(ctx, existingAccount.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
