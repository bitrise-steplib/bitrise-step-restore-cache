package step

import (
	"fmt"
	"strings"
	"time"

	"github.com/bitrise-steplib/bitrise-step-restore-cache/decompression"

	"github.com/bitrise-io/go-steputils/v2/cache/keytemplate"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

type Input struct {
	Verbose bool   `env:"verbose,required"`
	Key     string `env:"key,required"`
}

type Config struct {
	Verbose bool
	Keys    []string
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

	return &Config{
		Verbose: input.Verbose,
		Keys:    keySlice,
	}, nil
}

func (step RestoreCacheStep) Run(config *Config) {
	for _, key := range config.Keys {
		step.logger.Println()
		step.logger.Printf("Evaluating key template: %s", key)
		evaluatedKey, err := step.evaluateKey(key)
		if err != nil {
			step.logger.Warnf("Failed to evaluate key template: %s", key)
			continue
		}
		step.logger.Donef("Cache key: %s", evaluatedKey)

		step.logger.Println()
		step.logger.Printf("Restoring cache archive...")
		startTime := time.Now()
		if err := decompression.Decompress(evaluatedKey, step.logger, step.envRepo); err != nil {
			step.logger.Warnf("Failed to decompress cache archive: %s", evaluatedKey)
			continue
		}
		step.logger.Donef("Restored cache archive in %s", time.Since(startTime).Round(time.Second))
	}
}

func (step RestoreCacheStep) evaluateKey(keyTemplate string) (string, error) {
	model := keytemplate.NewModel(step.envRepo, step.logger)
	buildContext := keytemplate.BuildContext{
		Workflow:   step.envRepo.Get("BITRISE_WORKFLOW_ID"),
		Branch:     step.envRepo.Get("BITRISE_GIT_BRANCH"),
		CommitHash: step.envRepo.Get("BITRISE_GIT_COMMIT"),
	}
	return model.Evaluate(keyTemplate, buildContext)
}

// This method is currently here for debugging, but isn't used in the code, so golint complains.
// It may be removed in the future, but for now we'll leave it.
//nolint:golint,unused
func (step RestoreCacheStep) getArchiveContents(archivePath string) ([]string, error) {
	getArchiveContentsArgs := []string{
		"--list",
		"--file",
		archivePath,
	}

	cmd := step.commandFactory.Create("tar", getArchiveContentsArgs, nil)
	step.logger.Debugf("$ %s", cmd.PrintableCommandArgs())

	archiveContents, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		step.logger.Errorf("Failed to get archiveContents: %s", archiveContents)
		return nil, err
	}

	archiveContentsSlice := strings.Split(archiveContents, "\n")

	return archiveContentsSlice, nil
}
