package step

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestoreCacheStep_Run(t *testing.T) {
	mockLogger := new(mocks.Logger)
	mockLogger.On("Debugf", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("EnableDebugLog", mock.Anything).Return()
	mockLogger.On("Println", mock.Anything).Return()
	mockLogger.On("Printf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Donef", mock.Anything, mock.Anything).Return()

	tmpPath := t.TempDir()
	tmpEnvmanEnvstorePath := filepath.Join(tmpPath, "envstore.yml")
	require.NoError(t, writeToFile(tmpEnvmanEnvstorePath, ""))
	os.Setenv("ENVMAN_ENVSTORE_PATH", tmpEnvmanEnvstorePath)

	envRepo := mockEnvRepository{
		envs: map[string]string{
			"verbose":                "true",
			"BITRISEIO_ABCS_API_URL": "https://abcs.services.bitrise.io",
			"BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN": "", // SPECIFY TOKEN HERE
			"key": "restore-cache-testfile-5gb", // SPECIFY KEY HERE
		},
	}
	inputParser := stepconf.NewInputParser(envRepo)
	commandFactory := command.NewFactory(envRepo)

	step := RestoreCacheStep{
		logger:         mockLogger,
		inputParser:    inputParser,
		commandFactory: commandFactory,
		envRepo:        envRepo,
	}
	require.NoError(t, step.Run())
}

// ----------------
// --- Utilites ---

func writeToFile(path string, value string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(value)
	if err != nil {
		return err
	}

	return nil
}

type mockEnvRepository struct {
	envs map[string]string
}

func (m mockEnvRepository) List() []string {
	var envs []string
	for k, v := range m.envs {
		envs = append(envs, k+"="+v)
	}
	return envs
}

func (m mockEnvRepository) Unset(key string) error {
	delete(m.envs, key)
	return nil
}

func (m mockEnvRepository) Get(key string) string {
	return m.envs[key]
}

func (m mockEnvRepository) Set(key, value string) error {
	m.envs[key] = value
	return nil
}
