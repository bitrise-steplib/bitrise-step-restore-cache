package step

import (
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

type Input struct {
	Verbose bool `env:"verbose,required"`
}

type Config struct {
	Verbose bool
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

	return &Config{
		Verbose: input.Verbose,
	}, nil
}

func (step RestoreCacheStep) Run(config *Config) error {
	return nil
}
