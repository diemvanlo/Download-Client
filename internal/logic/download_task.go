package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"goload/internal/dataaccess/database"
	"goload/internal/dataaccess/file"
	"goload/internal/dataaccess/mq/producer"
	go_load "goload/internal/generated/downloadClient/v1"
	"goload/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type CreateDownloadTaskParams struct {
	Token        string
	DownloadType go_load.DownloadType
	URL          string
}

type CreateDownloadTaskOutput struct {
	DownloadTask *go_load.DownloadTask
}

type GetDownloadTaskListParams struct {
	Token  string
	Offset uint64
	Limit  uint64
}

type GetDownloadTaskListOutput struct {
	DownloadTaskList       []*go_load.DownloadTask
	TotalDownloadTaskCount uint64
}

type UpdateDownloadTaskParams struct {
	Token          string
	DownloadTaskID uint64
	URL            string
}

type UpdateDownloadTaskOutput struct {
	DownloadTask *go_load.DownloadTask
}

type DeleteDownloadTaskParams struct {
	Token          string
	DownloadTaskID uint64
}

type DownloadTask interface {
	CreateDownloadTask(context.Context, CreateDownloadTaskParams) (CreateDownloadTaskOutput, error)
	GetDownloadTaskList(context.Context, GetDownloadTaskListParams) (GetDownloadTaskListOutput, error)
	UpdateDownloadTask(context.Context, UpdateDownloadTaskParams) (UpdateDownloadTaskOutput, error)
	DeleteDownloadTask(context.Context, DeleteDownloadTaskParams) error
	ExecuteDownloadTask(context.Context, uint64) error
}

type downloadTask struct {
	tokenLogic                  Token
	accountDataAccessor         database.AccountDataAccessor
	downloadTaskDataAccessor    database.DownloadTaskDataAccessor
	downloadTaskCreatedProducer producer.DownloadTaskCreatedProducer
	goquDatabase                *goqu.Database
	logger                      *zap.Logger
	client                      file.Client
}

func NewDownloadTask(
	tokenLogic Token,
	accountDataAccessor database.AccountDataAccessor,
	downloadTaskDataAccessor database.DownloadTaskDataAccessor,
	downloadTaskCreatedProducer producer.DownloadTaskCreatedProducer,
	goquDatabase *goqu.Database,
	client file.Client,
	logger *zap.Logger,
) DownloadTask {
	return &downloadTask{
		tokenLogic:                  tokenLogic,
		accountDataAccessor:         accountDataAccessor,
		downloadTaskDataAccessor:    downloadTaskDataAccessor,
		downloadTaskCreatedProducer: downloadTaskCreatedProducer,
		goquDatabase:                goquDatabase,
		logger:                      logger,
	}
}

func (d downloadTask) databaseDownloadTaskToProtoDownloadTask(task database.DownloadTask, account database.Account) *go_load.DownloadTask {
	return &go_load.DownloadTask{
		Id: task.ID,
		OfAccount: &go_load.Account{
			Id:          account.ID,
			AccountName: account.AccountName,
		},
		DownloadType:   task.DownloadType,
		Url:            task.URL,
		DownloadStatus: go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING,
	}
}

func (d downloadTask) CreateDownloadTask(
	ctx context.Context,
	params CreateDownloadTaskParams,
) (CreateDownloadTaskOutput, error) {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return CreateDownloadTaskOutput{}, err
	}

	account, err := d.accountDataAccessor.GetAccountByID(ctx, accountID)
	if err != nil {
		return CreateDownloadTaskOutput{}, err
	}

	downloadTask := database.DownloadTask{
		OfAccountID:    accountID,
		DownloadType:   params.DownloadType,
		URL:            params.URL,
		DownloadStatus: go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING,
		Metadata: database.JSON{
			Data: make(map[string]any),
		},
	}

	txErr := d.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		downloadTaskID, createDownloadTaskErr := d.downloadTaskDataAccessor.
			WithDatabase(td).
			CreateDownloadTask(ctx, downloadTask)
		if createDownloadTaskErr != nil {
			return createDownloadTaskErr
		}

		downloadTask.ID = downloadTaskID
		produceErr := d.downloadTaskCreatedProducer.Produce(ctx, producer.DownloadTaskCreated{
			ID: downloadTaskID,
		})
		if produceErr != nil {
			return produceErr
		}

		return nil
	})
	if txErr != nil {
		return CreateDownloadTaskOutput{}, txErr
	}

	return CreateDownloadTaskOutput{
		DownloadTask: d.databaseDownloadTaskToProtoDownloadTask(downloadTask, account),
	}, nil
}

func (d downloadTask) GetDownloadTaskList(
	ctx context.Context,
	params GetDownloadTaskListParams,
) (GetDownloadTaskListOutput, error) {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return GetDownloadTaskListOutput{}, err
	}

	account, err := d.accountDataAccessor.GetAccountByID(ctx, accountID)
	if err != nil {
		return GetDownloadTaskListOutput{}, err
	}

	totalDownloadTaskCount, err := d.downloadTaskDataAccessor.GetDownloadTaskCountOfAccount(ctx, accountID)
	if err != nil {
		return GetDownloadTaskListOutput{}, err
	}

	downloadTaskList, err := d.downloadTaskDataAccessor.GetDownloadTaskListOfAccount(ctx, accountID, params.Offset, params.Limit)
	if err != nil {
		return GetDownloadTaskListOutput{}, err
	}

	return GetDownloadTaskListOutput{
		DownloadTaskList: lo.Map(downloadTaskList, func(task database.DownloadTask, index int) *go_load.DownloadTask {
			return d.databaseDownloadTaskToProtoDownloadTask(task, account)
		}),
		TotalDownloadTaskCount: totalDownloadTaskCount,
	}, nil
}

func (d downloadTask) UpdateDownloadTask(ctx context.Context, params UpdateDownloadTaskParams) (UpdateDownloadTaskOutput, error) {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return UpdateDownloadTaskOutput{}, err
	}

	account, err := d.accountDataAccessor.GetAccountByID(ctx, accountID)
	if err != nil {
		return UpdateDownloadTaskOutput{}, err
	}

	output := UpdateDownloadTaskOutput{}
	txErr := d.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		downloadTask, getDownloadTaskWithXLockErr := d.downloadTaskDataAccessor.
			WithDatabase(td).
			GetDownloadTaskWithXLock(ctx, params.DownloadTaskID)

		if getDownloadTaskWithXLockErr != nil {
			return getDownloadTaskWithXLockErr
		}

		if downloadTask.OfAccountID != accountID {
			return status.Error(codes.PermissionDenied, "permission denied")
		}

		downloadTask.URL = params.URL
		output.DownloadTask = d.databaseDownloadTaskToProtoDownloadTask(downloadTask, account)
		return d.downloadTaskDataAccessor.WithDatabase(td).UpdateDownloadTask(ctx, downloadTask)
	})

	if txErr != nil {
		return UpdateDownloadTaskOutput{}, txErr
	}

	return output, nil
}

func (d downloadTask) DeleteDownloadTask(ctx context.Context, params DeleteDownloadTaskParams) error {
	accountID, _, err := d.tokenLogic.GetAccountIDAndExpireTime(ctx, params.Token)
	if err != nil {
		return err
	}

	txErr := d.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		downloadTask, getDownloadTaskWithXLockErr := d.downloadTaskDataAccessor.
			WithDatabase(td).
			GetDownloadTaskWithXLock(ctx, params.DownloadTaskID)

		if getDownloadTaskWithXLockErr != nil {
			return getDownloadTaskWithXLockErr
		}

		if downloadTask.OfAccountID != accountID {
			return status.Error(codes.PermissionDenied, "permission denied")
		}

		return d.downloadTaskDataAccessor.WithDatabase(td).DeleteDownloadTask(ctx, params.DownloadTaskID)
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (d downloadTask) updateDownloadTaskStatusFromPendingToDownloading(
	ctx context.Context,
	id uint64) (bool, database.DownloadTask, error) {
	var (
		logger  = utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))
		updated = false
		task    database.DownloadTask
		err     error
	)

	txErr := d.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		downloadTask, err := d.downloadTaskDataAccessor.WithDatabase(td).GetDownloadTaskWithXLock(ctx, id)
		if err != nil {
			if errors.Is(err, database.ErrDownloadTaskNotFound) {
				logger.Warn("cannot find download task by id")
				return nil
			}

			logger.With(zap.Error(err)).Error("failed to get download task by id")
			return err
		}

		if downloadTask.DownloadStatus != go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING {
			logger.Warn("download task is not pending, will not execute")
			return nil
		}
		downloadTask.DownloadStatus = go_load.DownloadStatus_DOWNLOAD_STATUS_DOWNLOADING

		err = d.downloadTaskDataAccessor.UpdateDownloadTask(ctx, downloadTask)
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to update download task")
			return err
		}
		updated = true
		return nil
	})

	if txErr != nil {

		return false, database.DownloadTask{}, txErr
	}
	return updated, task, err
}

func (d downloadTask) ExecuteDownloadTask(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	updated, task, err := d.updateDownloadTaskStatusFromPendingToDownloading(ctx, id)
	if err != nil {
		return err
	}

	if !updated {
		return nil
	}

	var downloader Downloader
	switch task.DownloadType {
	case go_load.DownloadType_DOWNLOAD_TYPE_HTTP:
		downloader = NewHTTPDownloader(task.URL, d.logger)
	default:
		logger.With(zap.String("download_type", task.DownloadType.String())).Error("unsupported download type")
		return nil
	}

	fileName := fmt.Sprintf("download_file_%d", task.ID)
	fileWriterClose, err := d.client.Write(ctx, fileName)
	if err != nil {
		return err
	}

	defer fileWriterClose.Close()

	err = downloader.Download(ctx, fileWriterClose)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to download")
		return err
	}

	logger.With(zap.String("url", task.URL)).Info("downloaded")
	return nil
}
