package step

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/cache"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

type Input struct {
	Verbose bool   `env:"verbose,required"`
	Key     string `env:"key,required"`
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

func (step RestoreCacheStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return err
	}
	stepconf.Print(input)

	if strings.TrimSpace(input.Key) == "" {
		return fmt.Errorf("required input 'key' is empty")
	}

	step.logger.EnableDebugLog(input.Verbose)

	return cache.NewRestorer(step.envRepo, step.logger).Restore(cache.RestoreCacheInput{
		StepId:  "restore-cache",
		Verbose: input.Verbose,
		Keys:    strings.Split(input.Key, "\n"),
	})
}
