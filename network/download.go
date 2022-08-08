package network

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/retryhttp"
)

type DownloadParams struct {
	APIBaseURL   string
	Token        string
	CacheKeys    []string
	DownloadPath string
}

var ErrCacheNotFound = errors.New("no cache archive found for the provided keys")

func Download(params DownloadParams, logger log.Logger) error {
	client := newApiClient(retryhttp.NewClient(logger), params.APIBaseURL, params.Token)

	logger.Debugf("Get download URL")
	url, err := client.restore(params.CacheKeys)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	logger.Debugf("Download archive")
	file, err := os.Create(params.DownloadPath)
	if err != nil {
		return fmt.Errorf("can't open download location: %w", err)
	}
	defer file.Close()

	respBody, err := client.downloadArchive(url)
	if err != nil {
		return fmt.Errorf("failed to download archive: %w", err)
	}
	defer respBody.Close()
	_, err = io.Copy(file, respBody)
	if err != nil {
		return fmt.Errorf("failed to save archive to disk: %w", err)
	}

	return nil
}
