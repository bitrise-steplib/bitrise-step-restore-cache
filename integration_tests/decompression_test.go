//go:build integration

package integration_tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-steplib/bitrise-step-restore-cache/decompression"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/stretchr/testify/assert"
)

func Test_Decompression(t *testing.T) {
	checkForZSTDBinary()

	tests := []struct {
		name        string
		archivePath string
		wantErr     bool
	}{
		{
			name:        "Single Item Archive",
			archivePath: "testdata/single-item.tzst",
			wantErr:     false,
		},
		{
			name:        "Single Directory Archive",
			archivePath: "testdata/single-directory.tzst",
			wantErr:     false,
		},
		{
			name:        "Multiple Item Archive",
			archivePath: "testdata/multiple-items.tzst",
			wantErr:     false,
		},
		{
			name:        "Nonexistent Archive",
			archivePath: "testdata/nonexistent.tzst",
			wantErr:     true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			logger := log.NewLogger()
			envRepo := env.NewRepository()
			tempDir := t.TempDir()

			// When
			decompressionErr := decompression.Decompress(
				testCase.archivePath,
				logger,
				envRepo,
				"--directory", tempDir,
			)

			// Then
			if testCase.wantErr {
				assert.Error(t, decompressionErr)
				return
			} else {
				assert.NoError(t, decompressionErr)
			}

			expectedArchiveContents, err := listArchiveContents(testCase.archivePath)
			if err != nil {
				t.Errorf("Failed to list archive contents: %v", err)
			}

			var actualDecompressedContents []string
			if err = filepath.Walk(
				tempDir,
				func(path string, info os.FileInfo, err error) error {
					// This walks the temp directory, and converts the paths to relative paths
					// to match the output of the tar command used in `listArchiveContents`.
					if err != nil {
						return err
					}
					if path == tempDir {
						return nil
					}
					if info.IsDir() {
						path = path + string(os.PathSeparator)
					}
					path = strings.TrimPrefix(path, tempDir)
					path = strings.TrimPrefix(path, string(os.PathSeparator))
					if len(path) > 0 {
						actualDecompressedContents = append(actualDecompressedContents, path)
					}
					return nil
				},
			); err != nil {
				t.Errorf("Failed to walk temp dir: %v", err)
			}

			assert.NoError(t, err)
			assert.ElementsMatch(t, actualDecompressedContents, expectedArchiveContents)
		})
	}
}

func checkForZSTDBinary() {
	_, err := exec.LookPath("zstd")
	if err != nil {
		panic("zstd is required for integration tests")
	}
}

func listArchiveContents(archivePath string) ([]string, error) {
	listArchiveContentsArgs := []string{
		"--list",
		"--file",
		archivePath,
	}

	commandFactory := command.NewFactory(env.NewRepository())
	cmd := commandFactory.Create("tar", listArchiveContentsArgs, nil)

	archiveContents, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, err
	}

	archiveContentsSlice := strings.Split(archiveContents, "\n")

	return archiveContentsSlice, nil
}
