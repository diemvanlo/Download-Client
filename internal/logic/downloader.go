package logic

import (
	"context"
	"go.uber.org/zap"
	"goload/internal/utils"
	"io"
	"net/http"
)

type Downloader interface {
	Download(ctx context.Context, writer io.Writer) error
}

type HTTPDownloader struct {
	URL    string
	logger *zap.Logger
}

func NewHTTPDownloader(URL string, logger *zap.Logger) Downloader {
	return &HTTPDownloader{
		URL:    URL,
		logger: logger,
	}
}

func (h *HTTPDownloader) Download(ctx context.Context, writer io.Writer) error {
	logger := utils.LoggerWithContext(ctx, h.logger)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, h.URL, http.NoBody)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create http get request")
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to make http get request")
	}

	defer response.Body.Close()

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to copy response body")
		return err
	}

	return nil
}
