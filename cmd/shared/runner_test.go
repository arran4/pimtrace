package shared

import (
	"bytes"
	"errors"
	"fmt"
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
			wantStdout: "Usage:  testtool [Flags] [Query]",
			wantStderr: "",
		},
		{
			name:       "short help flag",
			args:       []string{"-h"},
			wantErr:    0,
			wantStdout: "Usage:  testtool [Flags] [Query]",
			wantStderr: "",
		},
		{
			name:       "version flag",
			args:       []string{"-version"},
			wantErr:    0,
			wantStdout: "v1.0.0",
			wantStderr: "",
			setupCfg: func(cfg *Config) {
				cfg.PrintVersion = func(w io.Writer) {
					_, _ = fmt.Fprintln(w, "v1.0.0")
				}
			},
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
			wantStdout: "success output", // Usually no stdout unless PrintQueryHelp or output handler writes
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
		{
			name:       "malformed query",
			args:       []string{"-parser", "basic", "filter", "invalid_token"}, // missing arguments to filter
			wantErr:    2,
			wantStdout: "",
			wantStderr: "Parse Error:",
			setupCfg: func(cfg *Config) {
				cfg.InputHandler = func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
					return &mockData{}, nil
				}
			},
		},
		{
			name:       "exact trailing tokens",
			args:       []string{"-parser", "basic", "filter", "c.name", "eq", ".a b c"}, // .a b c is passed as one token
			wantErr:    0,
			wantStdout: "success output",
			wantStderr: "",
			setupCfg: func(cfg *Config) {
				cfg.InputHandler = func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
					return &mockData{}, nil
				}
				cfg.OutputHandler = func(data pimtrace.Data, outputType string, outputFile string, stdout io.Writer) error {
					_, _ = stdout.Write([]byte("success output"))
					return nil
				}
			},
		},
		{
			name:       "output write failure",
			args:       []string{"-parser", "basic", "filter", ".1", "eq", ".1"},
			wantErr:    1,
			wantStdout: "",
			wantStderr: "Write Error: mock write error",
			setupCfg: func(cfg *Config) {
				cfg.InputHandler = func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
					return &mockData{}, nil
				}
				cfg.OutputHandler = func(data pimtrace.Data, outputType string, outputFile string, stdout io.Writer) error {
					return errors.New("mock write error")
				}
			},
		},
		{
			name:       "injected stdin reader regression",
			args:       []string{"-parser", "basic", "filter", ".1", "eq", ".1"},
			wantErr:    0,
			wantStdout: "",
			wantStderr: "",
			setupCfg: func(cfg *Config) {
				cfg.Stdin = strings.NewReader("test input")
				cfg.InputHandler = func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
					var foundReader io.Reader
					var foundWriter io.Writer

					for _, op := range ops {
						if r, ok := op.(io.Reader); ok {
							foundReader = r
						}
						if w, ok := op.(io.Writer); ok {
							foundWriter = w
						}
					}

					if foundReader == nil {
						return nil, errors.New("injected Stdin reader not found")
					}

					if foundWriter != nil {
						// Since we wrapped c.Stdin in an anonymous struct to ONLY implement io.Reader,
						// it shouldn't also match as an io.Writer.
						return nil, errors.New("injected Stdin reader incorrectly cast to io.Writer")
					}

					buf, err := io.ReadAll(foundReader)
					if err != nil {
						return nil, fmt.Errorf("failed to read from injected reader: %w", err)
					}

					if string(buf) != "test input" {
						return nil, fmt.Errorf("read unexpected input from injected reader: %q", string(buf))
					}

					return &mockData{}, nil
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
			} else if tt.wantStderr == "" && len(strings.TrimSpace(stderr.String())) > 0 {
				t.Errorf("Run() stderr expected empty, got %q", stderr.String())
			}

			if tt.wantStdout != "" && !strings.Contains(stdout.String(), strings.TrimSpace(tt.wantStdout)) {
				t.Errorf("Run() stdout = %q, want containing %q", stdout.String(), tt.wantStdout)
			} else if tt.wantStdout == "" && len(strings.TrimSpace(stdout.String())) > 0 {
				t.Errorf("Run() stdout expected empty, got %q", stdout.String())
			}
		})
	}
}
