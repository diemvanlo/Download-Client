package database

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"goload/internal/generated/grpc/go_load"
	"goload/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	TabNameDownloadTasks = goqu.T("download_tasks")
)

const (
	ColNameDownloadTaskId             = "id"
	ColNameDownloadTaskOfAccountID    = "of_account_id"
	ColNameDownloadTaskDownloadType   = "download_type"
	ColNameDownloadTaskURL            = "url"
	ColNameDownloadTaskDownloadStatus = "download_status"
	ColNameDownloadTaskMetadata       = "metadata"
)

type DownloadTask struct {
	ID             uint64                 `db:"id" goqu:"skipinsert,skipupdate"`
	OfAccountID    uint64                 `db:"of_account_id"`
	DownloadType   go_load.DownloadType   `db:"download_type"`
	URL            string                 `db:"url"`
	DownloadStatus go_load.DownloadStatus `db:"download_status"`
	Metadata       JSON                   `db:"metadata"`
}

type DownloadTaskDataAccessor interface {
	CreateDownloadTask(ctx context.Context, task DownloadTask) (uint64, error)
	GetDownloadTaskListOfAccount(ctx context.Context, accountID, offset, limit uint64) ([]DownloadTask, error)
	GetDownloadTaskCountOfAccount(ctx context.Context, accountId uint64) (uint64, error)
	GetDownloadTask(ctx context.Context, id uint64) (DownloadTask, error)
	GetDownloadTaskWithXLock(ctx context.Context, id uint64) (DownloadTask, error)
	UpdateDownloadTask(ctx context.Context, task DownloadTask) error
	DeleteDownloadTask(ctx context.Context, id uint64) error
	WithDatabase(database Database) DownloadTaskDataAccessor
}

type downloadTaskDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewDownloadTaskDataAccessor(
	database *goqu.Database,
	logger *zap.Logger,
) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (d downloadTaskDataAccessor) CreateDownloadTask(ctx context.Context, task DownloadTask) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("task", task))

	result, err := d.database.
		Insert(TabNameDownloadTasks).
		Rows(task).
		Executor().
		ExecContext(ctx)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create download task")
		return 0, status.Error(codes.Internal, "failed to create download task")
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get last inserted id")
		return 0, status.Error(codes.Internal, "failed to get last inserted id")
	}

	return uint64(lastInsertedID), nil
}

func (d downloadTaskDataAccessor) DeleteDownloadTask(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	if _, err := d.database.
		Delete(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameDownloadTaskId: id}).
		Executor().
		ExecContext(ctx); err != nil {
		logger.With(zap.Error(err)).Error("failed to delete download task")
		return status.Error(codes.Internal, "failed to delete download task")
	}

	return nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskCountOfAccount(ctx context.Context, userID uint64) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("account_id", userID))

	count, err := d.database.
		From(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameDownloadTaskOfAccountID: userID}).
		CountContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task count of account")
		return 0, status.Error(codes.Internal, "failed to get download task count of account")
	}

	return uint64(count), nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskListOfAccount(ctx context.Context, userID uint64, offset uint64, limit uint64) ([]DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("account_id", userID))

	tasks := []DownloadTask{}
	err := d.database.
		From(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameDownloadTaskOfAccountID: userID}).
		Offset(uint(offset)).
		Limit(uint(limit)).
		ScanStructsContext(ctx, &tasks)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task list of account")
		return nil, status.Error(codes.Internal, "failed to get download task list of account")
	}

	return tasks, nil
}

func (d downloadTaskDataAccessor) GetDownloadTask(ctx context.Context, id uint64) (DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	task := DownloadTask{}
	found, err := d.database.
		From(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameDownloadTaskId: id}).
		ScanStructContext(ctx, &task)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task by id")
		return DownloadTask{}, status.Error(codes.Internal, "failed to get download task by id")
	}

	if !found {
		logger.Warn("cannot find download task by id")
		return DownloadTask{}, status.Error(codes.NotFound, "cannot find download task by id")
	}

	return task, nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskWithXLock(ctx context.Context, id uint64) (DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	task := DownloadTask{}
	found, err := d.database.
		From(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameDownloadTaskId: id}).
		ForUpdate(goqu.Wait).
		ScanStructContext(ctx, &task)

	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task by id")
		return DownloadTask{}, status.Error(codes.Internal, "failed to get download task by id")
	}

	if !found {
		logger.Warn("cannot find download task by id")
		return DownloadTask{}, status.Error(codes.NotFound, "cannot find download task by id")
	}

	return task, nil
}

func (d downloadTaskDataAccessor) UpdateDownloadTask(ctx context.Context, task DownloadTask) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("task", task))

	if _, err := d.database.
		Update(TabNameDownloadTasks).
		Set(task).
		Where(goqu.Ex{ColNameDownloadTaskId: task.ID}).
		Executor().
		ExecContext(ctx); err != nil {
		logger.With(zap.Error(err)).Error("failed to update download task")
		return status.Errorf(codes.Internal, "failed to update download task")
	}

	return nil
}

func (d downloadTaskDataAccessor) WithDatabase(database Database) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
		logger:   d.logger,
	}
}
