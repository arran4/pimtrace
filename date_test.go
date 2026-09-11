package pimtrace_test

import (
	"pimtrace"
	"testing"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"2020-01-01", false},
		{"2020-01-01T00:00:00Z", false},
		{"2020-01-02 00:00:00", false},
		{"20231027T100000Z", false},
		{"20231027", false},
		{"Mon, 02 Jan 2006 15:04:05 -0700", false},
		{"", true},
		{"   ", true},
		{"apple", true},
		{"invalid date", true},
		{"not-a-date", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := pimtrace.ParseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ParseDate(%q) returned nil time without error", tt.input)
			}
		})
	}
}

func TestCoerceDate(t *testing.T) {
	tests := []struct {
		name    string
		input   pimtrace.Value
		wantErr bool
	}{
		{"valid string date", pimtrace.SimpleStringValue("2020-01-01"), false},
		{"valid ical date", pimtrace.SimpleStringValue("20231027T100000Z"), false},
		{"integer is not date", pimtrace.SimpleIntegerValue(42), true},
		{"float is not date", pimtrace.SimpleFloatValue(10.25), true},
		{"empty string is error", pimtrace.SimpleStringValue(""), true},
		{"invalid text is error", pimtrace.SimpleStringValue("apple"), true},
		{"nil value is error", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pimtrace.CoerceDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoerceDate(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("CoerceDate(%v) returned nil time without error", tt.input)
			}
		})
	}
}
