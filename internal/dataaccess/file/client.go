package file

import (
	"bufio"
	"context"
	"go.uber.org/zap"
	"goload/internal/configs"
	"goload/internal/utils"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"path"
)

type Client interface {
	Write(ctx context.Context, filePath string) (io.WriteCloser, error)
	Read(ctx context.Context, filePath string) (io.ReadCloser, error)
}

type bufferFileReader struct {
	file           *os.File
	bufferedReader *bufio.Reader
}

func newBufferFileReader(file *os.File) io.ReadCloser {
	return &bufferFileReader{
		file:           file,
		bufferedReader: bufio.NewReader(file),
	}
}

func (b bufferFileReader) Read(p []byte) (int, error) {
	return b.bufferedReader.Read(p)
}

func (b bufferFileReader) Close() error {
	return b.file.Close()
}

type localFileClient struct {
	logger            *zap.Logger
	downloadDirectory string
}

func NewLocalFileClient(logger *zap.Logger, config configs.DownloadConfig) Client {
	return &localFileClient{
		logger:            logger,
		downloadDirectory: config.DownloadDir,
	}
}

func (l localFileClient) Write(ctx context.Context, filePath string) (io.WriteCloser, error) {
	logger := utils.LoggerWithContext(ctx, l.logger).With(zap.String("filePath", filePath))

	absolutFilePath := path.Join(l.downloadDirectory, filePath)
	file, err := os.Open(absolutFilePath)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to open file")
		return nil, status.Errorf(status.Code(err), "failed to open file: %v", err)
	}

	return file, nil
}

func (l localFileClient) Read(ctx context.Context, filePath string) (io.ReadCloser, error) {
	logger := utils.LoggerWithContext(ctx, l.logger).With(zap.String("filePath", filePath))
	absolutFilePath := path.Join(l.downloadDirectory, filePath)
	file, err := os.Open(absolutFilePath)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to open file")
		return nil, status.Errorf(status.Code(err), "failed to open file: %v", err)
	}

	return newBufferFileReader(file), nil
}
