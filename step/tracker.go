package step

import (
	"io/fs"
	"time"

	"github.com/bitrise-io/go-utils/v2/analytics"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

type stepTracker struct {
	tracker analytics.Tracker
	logger  log.Logger
}

func newStepTracker(config Config, envRepo env.Repository, logger log.Logger) stepTracker {
	p := analytics.Properties{
		"is_pr_build": envRepo.Get("IS_PR") == "true",
		"key_count":   len(config.Keys),
	}
	return stepTracker{
		tracker: analytics.NewDefaultTracker(logger, p),
		logger:  logger,
	}
}

func (t *stepTracker) logArchiveDownloaded(downloadTime time.Duration, info fs.FileInfo) {
	properties := analytics.Properties{
		"download_time_s":     downloadTime.Truncate(time.Second).Seconds(),
		"download_size_bytes": info.Size(),
	}
	t.tracker.Enqueue("step_restore_cache_archive_downloaded", properties)
}

func (t *stepTracker) logArchiveExtracted(extractionTime time.Duration) {
	properties := analytics.Properties{
		"extraction_time_s": extractionTime.Truncate(time.Second).Seconds(),
	}
	t.tracker.Enqueue("step_restore_cache_archive_extracted", properties)
}

func (t *stepTracker) wait() {
	t.tracker.Wait()
}
