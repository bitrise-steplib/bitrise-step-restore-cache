package main

import (
	"os"

	"github.com/bitrise-steplib/bitrise-step-restore-cache/step"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()
	envRepo := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepo)

	restoreCacheStep := step.New(logger, inputParser, envRepo)

	exitCode := 0

	config, err := restoreCacheStep.ProcessConfig()
	if err != nil {
		logger.Errorf(err.Error())
		exitCode = 1
		return exitCode
	}

	logger.EnableDebugLog(config.Verbose)

	if err := restoreCacheStep.Run(config); err != nil {
		logger.Errorf(err.Error())
		exitCode = 1
		return exitCode
	}

	return exitCode
}
