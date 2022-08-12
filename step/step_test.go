package step

import (
	"reflect"
	"testing"

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
				Verbose:        true,
				Keys:           []string{"valid-key"},
				APIBaseURL:     "fake service URL",
				APIAccessToken: "fake access token",
			},
			wantErr: false,
		},
		{
			name: "Valid key input with multiple keys",
			inputParser: fakeInputParser{
				verbose: true,
				key:     "valid-key\nvalid-key-2",
			},
			want: &Config{
				Verbose:        true,
				Keys:           []string{"valid-key", "valid-key-2"},
				APIBaseURL:     "fake service URL",
				APIAccessToken: "fake access token",
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
				envRepo: fakeEnvRepo{envVars: map[string]string{
					"BITRISEIO_CACHE_SERVICE_URL":          "fake service URL",
					"BITRISEIO_CACHE_SERVICE_ACCESS_TOKEN": "fake access token",
				}},
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

func Test_evaluateKeys(t *testing.T) {
	type args struct {
		keys    []string
		envRepo fakeEnvRepo
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Happy path",
			args: args{
				keys: []string{"npm-cache-{{ .Branch }}"},
				envRepo: fakeEnvRepo{
					envVars: map[string]string{
						"BITRISE_WORKFLOW_ID": "primary",
						"BITRISE_GIT_BRANCH":  "main",
						"BITRISE_GIT_COMMIT":  "9de033412f24b70b59ca8392ccb9f61ac5af4cc3",
					},
				},
			},
			want:    []string{"npm-cache-main"},
			wantErr: false,
		},
		{
			name: "Multiple keys",
			args: args{
				keys: []string{
					"npm-cache-{{ .Branch }}",
					"npm-cache-",
					"",
				},
				envRepo: fakeEnvRepo{
					envVars: map[string]string{
						"BITRISE_WORKFLOW_ID": "primary",
						"BITRISE_GIT_BRANCH":  "main",
						"BITRISE_GIT_COMMIT":  "9de033412f24b70b59ca8392ccb9f61ac5af4cc3",
					},
				},
			},
			want: []string{
				"npm-cache-main",
				"npm-cache-",
			},
			wantErr: false,
		},
		{
			name: "Empty environment variables",
			args: args{
				keys:    []string{"npm-cache-{{ .Branch }}"},
				envRepo: fakeEnvRepo{},
			},
			want:    []string{"npm-cache-"},
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
			evaluatedKeys, err := step.evaluateKeys(testCase.args.keys)
			if (err != nil) != testCase.wantErr {
				t.Errorf("evaluateKey() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if !reflect.DeepEqual(evaluatedKeys, testCase.want) {
				t.Errorf("evaluateKey() = %v, want %v", evaluatedKeys, testCase.want)
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
