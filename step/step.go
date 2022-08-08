package step

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-steputils/v2/cache/keytemplate"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-restore-cache/decompression"
	"github.com/bitrise-steplib/bitrise-step-restore-cache/network"
)

type Input struct {
	Verbose bool   `env:"verbose,required"`
	Key     string `env:"key,required"`
}

type Config struct {
	Verbose        bool
	Keys           []string
	APIBaseURL     stepconf.Secret
	APIAccessToken stepconf.Secret
}

type RestoreCacheStep struct {
	logger         log.Logger
	inputParser    stepconf.InputParser
	commandFactory command.Factory
	envRepo        env.Repository
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	commandFactory command.Factory,
	envRepo env.Repository,
) RestoreCacheStep {
	return RestoreCacheStep{
		logger:         logger,
		inputParser:    inputParser,
		commandFactory: commandFactory,
		envRepo:        envRepo,
	}
}

func (step RestoreCacheStep) ProcessConfig() (*Config, error) {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return nil, err
	}
	stepconf.Print(input)

	if strings.TrimSpace(input.Key) == "" {
		return nil, fmt.Errorf("required input 'key' is empty")
	}
	keySlice := strings.Split(input.Key, "\n")

	apiBaseURL := step.envRepo.Get("BITRISEIO_CACHE_SERVICE_URL")
	if apiBaseURL == "" {
		return nil, fmt.Errorf("the secret 'BITRISEIO_CACHE_SERVICE_URL' is not defined")
	}
	apiAccessToken := step.envRepo.Get("BITRISEIO_CACHE_SERVICE_ACCESS_TOKEN")
	if apiAccessToken == "" {
		return nil, fmt.Errorf("the secret 'BITRISEIO_CACHE_SERVICE_ACCESS_TOKEN' is not defined")
	}

	return &Config{
		Verbose:        input.Verbose,
		Keys:           keySlice,
		APIBaseURL:     stepconf.Secret(apiBaseURL),
		APIAccessToken: stepconf.Secret(apiAccessToken),
	}, nil
}

func (step RestoreCacheStep) Run(config *Config) error {
	var evaluatedKeys []string

	for _, key := range config.Keys {
		step.logger.Println()
		step.logger.Printf("Evaluating key template: %s", key)
		evaluatedKey, err := step.evaluateKey(key)
		if err != nil {
			return fmt.Errorf("failed to evaluate key template: %s", err)
		}
		step.logger.Donef("Cache key: %s", evaluatedKey)
		evaluatedKeys = append(evaluatedKeys, evaluatedKey)
	}

	step.logger.Println()
	step.logger.Infof("Downloading archive...")
	downloadStartTime := time.Now()
	archivePath, err := step.download(evaluatedKeys, *config)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	step.logger.Donef("Downloaded archive in %s", time.Since(downloadStartTime).Round(time.Second))

	step.logger.Println()
	step.logger.Infof("Restoring archive...")
	startTime := time.Now()
	if err := decompression.Decompress(archivePath, step.logger, step.envRepo); err != nil {
		return fmt.Errorf("failed to decompress cache archive: %w", err)
	}
	step.logger.Donef("Restored archive in %s", time.Since(startTime).Round(time.Second))

	return nil
}

func (step RestoreCacheStep) evaluateKey(keyTemplate string) (string, error) {
	model := keytemplate.NewModel(step.envRepo, step.logger)
	buildContext := keytemplate.BuildContext{
		Workflow:   step.envRepo.Get("BITRISE_TRIGGERED_WORKFLOW_ID"),
		Branch:     step.envRepo.Get("BITRISE_GIT_BRANCH"),
		CommitHash: step.envRepo.Get("BITRISE_GIT_COMMIT"),
	}
	return model.Evaluate(keyTemplate, buildContext)
}

func (step RestoreCacheStep) download(keys []string, config Config) (string, error) {
	dir, err := os.MkdirTemp("", "step-restore-cache")
	if err != nil {
		return "", err
	}
	name := fmt.Sprintf("cache-%s.tzst", time.Now().UTC().Format("20060102-150405"))
	downloadPath := filepath.Join(dir, name)

	params := network.DownloadParams{
		APIBaseURL:   string(config.APIBaseURL),
		Token:        string(config.APIAccessToken),
		CacheKeys:    keys,
		DownloadPath: downloadPath,
	}
	err = network.Download(params, step.logger)
	if err != nil {
		return "", err
	}

	step.logger.Debugf("Archive downloaded to %s", downloadPath)

	return downloadPath, nil
}
