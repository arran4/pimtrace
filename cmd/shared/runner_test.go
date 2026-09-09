package shared

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"pimtrace"
)

type mockData struct{}

func (m *mockData) Truncate(n int) pimtrace.Data                       { return m }
func (m *mockData) SetEntry(n int, entry pimtrace.Entry) pimtrace.Data { return m }
func (m *mockData) Len() int                                           { return 0 }
func (m *mockData) Entry(n int) pimtrace.Entry                         { return nil }
func (m *mockData) NewSelf() pimtrace.Data                             { return m }

func TestRunner(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantErr    int
		wantStdout string
		wantStderr string
		setupCfg   func(cfg *Config)
	}{
		{
			name:       "help flag",
			args:       []string{"-help"},
			wantErr:    0,
			wantStdout: "",
			wantStderr: "Usage:  testtool [Flags] [Query]\n",
		},
		{
			name:       "version flag",
			args:       []string{"-version"},
			wantErr:    0,
			wantStdout: "",
			wantStderr: "",
		},
		{
			name:       "no args",
			args:       []string{},
			wantErr:    2,
			wantStdout: "",
			wantStderr: "No query found\n",
		},
		{
			name:       "invalid flag",
			args:       []string{"-invalidflag"},
			wantErr:    2,
			wantStdout: "",
			wantStderr: "flag provided but not defined: -invalidflag\n",
		},
		{
			name:       "missing parser argument",
			args:       []string{"-parser", "invalid", "myquery"},
			wantErr:    2,
			wantStdout: "",
			wantStderr: "Please use -parser=basic parameter",
		},
		{
			name:       "input handler error",
			args:       []string{"-parser", "basic", "query"},
			wantErr:    1,
			wantStdout: "",
			wantStderr: "Read Error: mock read error",
			setupCfg: func(cfg *Config) {
				cfg.InputHandler = func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
					return nil, errors.New("mock read error")
				}
			},
		},
		{
			name:       "successful query",
			args:       []string{"-parser", "basic", "filter", ".1", "eq", ".1"},
			wantErr:    0,
			wantStdout: "", // Usually no stdout unless PrintQueryHelp or output handler writes
			wantStderr: "",
			setupCfg: func(cfg *Config) {
				cfg.InputHandler = func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
					// Need non-nil Data to avoid nil pointer panic in ast.Filter
					return &mockData{}, nil
				}
				cfg.OutputHandler = func(data pimtrace.Data, outputType string, outputFile string, stdout io.Writer) error {
					_, _ = stdout.Write([]byte("success output"))
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			cfg := &Config{
				Stdout: &stdout,
				Stderr: &stderr,
				Args:   tt.args,
				Name:   "testtool",
				InputHandler: func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
					return nil, nil
				},
				OutputHandler: func(data pimtrace.Data, outputType string, outputFile string, w io.Writer) error {
					return nil
				},
			}

			if tt.setupCfg != nil {
				tt.setupCfg(cfg)
			}

			if got := Run(cfg); got != tt.wantErr {
				t.Errorf("Run() = %v, want %v", got, tt.wantErr)
			}

			if tt.wantStderr != "" && !strings.Contains(stderr.String(), strings.TrimSpace(tt.wantStderr)) {
				t.Errorf("Run() stderr = %q, want containing %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}
