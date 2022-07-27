package step

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

func Test_ProcessConfig(t *testing.T) {
	tests := []struct {
		name        string
		inputParser fakeInputParser
		want        *Config
		wantErr     bool
	}{
		{
			name: "Valid key input",
			inputParser: fakeInputParser{
				verbose: true,
				key:     "valid-key",
			},
			want: &Config{
				Verbose: true,
				Key:     "valid-key",
			},
			wantErr: false,
		},
		{
			name: "Invalid key input",
			inputParser: fakeInputParser{
				verbose: false,
				key:     "  ",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			step := RestoreCacheStep{
				logger:      log.NewLogger(),
				inputParser: testCase.inputParser,
				envRepo:     env.NewRepository(),
			}

			// When
			processedConfig, err := step.ProcessConfig()

			// Then
			if (err != nil) != testCase.wantErr {
				t.Errorf("ProcessConfig() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if !reflect.DeepEqual(processedConfig, testCase.want) {
				t.Errorf("ProcessConfig() = %v, want %v", processedConfig, testCase.want)
			}
		})
	}
}

func Test_evaluateKey(t *testing.T) {
	type args struct {
		key     string
		envRepo fakeEnvRepo
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy path",
			args: args{
				key: "npm-cache-{{ .Branch }}",
				envRepo: fakeEnvRepo{
					envVars: map[string]string{
						"BITRISE_WORKFLOW_ID": "primary",
						"BITRISE_GIT_BRANCH":  "main",
						"BITRISE_GIT_COMMIT":  "9de033412f24b70b59ca8392ccb9f61ac5af4cc3",
					},
				},
			},
			want:    "npm-cache-main",
			wantErr: false,
		},
		{
			name: "Empty environment variables",
			args: args{
				key:     "npm-cache-{{ .Branch }}",
				envRepo: fakeEnvRepo{},
			},
			want:    "npm-cache-",
			wantErr: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			step := RestoreCacheStep{
				logger:      log.NewLogger(),
				inputParser: fakeInputParser{},
				envRepo:     testCase.args.envRepo,
			}

			// When
			evaluatedKey, err := step.evaluateKey(testCase.args.key)
			if (err != nil) != testCase.wantErr {
				t.Errorf("evaluateKey() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if evaluatedKey != testCase.want {
				t.Errorf("evaluateKey() = %v, want %v", evaluatedKey, testCase.want)
			}
		})
	}
}

func Test_getArchiveContents(t *testing.T) {
	tests := []struct {
		name        string
		archivePath string
		want        []string
		wantErr     bool
	}{
		{
			name:        "Single File Archive",
			archivePath: "testdata/single-file-archive.tzst",
			want:        []string{"dummy_file.txt"},
			wantErr:     false,
		},
		{
			name:        "Single Directory Archive",
			archivePath: "testdata/single-directory-archive.tzst",
			want:        []string{"subfolder/", "subfolder/nested_file.txt"},
			wantErr:     false,
		},
		{
			name:        "Multiple Item Archive",
			archivePath: "testdata/multi-item-archive.tzst",
			want:        []string{"subfolder/", "subfolder/nested_file.txt", "dummy_file.txt"},
			wantErr:     false,
		},
		{
			name:        "Nonexistent Archive",
			archivePath: "testdata/nonexistent-archive.tzst",
			want:        nil,
			wantErr:     true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			step := RestoreCacheStep{
				logger:         log.NewLogger(),
				commandFactory: command.NewFactory(env.NewRepository()),
			}

			// When
			archiveContents, err := step.getArchiveContents(testCase.archivePath)

			// Then
			if (err != nil) != testCase.wantErr {
				t.Errorf("getArchiveContents() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if !reflect.DeepEqual(archiveContents, testCase.want) {
				t.Errorf("getArchiveContents() = %v, want %v", archiveContents, testCase.want)
			}
		})
	}
}

// Helpers

type fakeInputParser struct {
	verbose bool
	key     string
}

func (f fakeInputParser) Parse(input interface{}) error {
	inputRef := input.(*Input)
	inputRef.Verbose = f.verbose
	inputRef.Key = f.key
	return nil
}

type fakeEnvRepo struct {
	envVars map[string]string
}

func (repo fakeEnvRepo) Get(key string) string {
	value, ok := repo.envVars[key]
	if ok {
		return value
	}
	return ""
}

func (repo fakeEnvRepo) Set(key, value string) error {
	repo.envVars[key] = value
	return nil
}

func (repo fakeEnvRepo) Unset(key string) error {
	repo.envVars[key] = ""
	return nil
}

func (repo fakeEnvRepo) List() (values []string) {
	for _, value := range repo.envVars {
		values = append(values, value)
	}
	return
}
