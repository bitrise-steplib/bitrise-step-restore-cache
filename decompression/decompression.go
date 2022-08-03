package decompression

import (
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

func Decompress(archivePath string, logger log.Logger, envRepo env.Repository, additionalArgs ...string) error {
	commandFactory := command.NewFactory(envRepo)

	decompressTarArgs := []string{
		"--use-compress-program",
		"zstd -d",
		"-xf",
		archivePath,
	}

	if len(additionalArgs) > 0 {
		decompressTarArgs = append(decompressTarArgs, additionalArgs...)
	}

	cmd := commandFactory.Create("tar", decompressTarArgs, nil)
	logger.Debugf("$ %s", cmd.PrintableCommandArgs())

	output, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		logger.Errorf("Failed to decompress cache archive: %s", output)
		return err
	}

	return nil
}
