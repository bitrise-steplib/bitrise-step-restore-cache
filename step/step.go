package step

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/cache/keytemplate"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

type Input struct {
	Verbose bool   `env:"verbose,required"`
	Key     string `env:"key,required"`
}

type Config struct {
	Verbose bool
	Key     string
}

type RestoreCacheStep struct {
	logger      log.Logger
	inputParser stepconf.InputParser
	envRepo     env.Repository
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	envRepo env.Repository,
) RestoreCacheStep {
	return RestoreCacheStep{
		logger:      logger,
		inputParser: inputParser,
		envRepo:     envRepo,
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

	return &Config{
		Verbose: input.Verbose,
		Key:     input.Key,
	}, nil
}

func (step RestoreCacheStep) Run(config *Config) error {
	evaluatedKey, err := step.evaluateKey(config.Key)
	if err != nil {
		return err
	}
	step.logger.Donef("Cache key: %s", evaluatedKey)
	return nil
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
