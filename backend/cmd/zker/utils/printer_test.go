package utils

import (
	"testing"
)

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		maxLen  int
		want    string
	}{
		{
			name:   "string shorter than max",
			s:      "short",
			maxLen: 10,
			want:   "short",
		},
		{
			name:   "string equal to max",
			s:      "exactlen",
			maxLen: 8,
			want:   "exactlen",
		},
		{
			name:   "string longer than max",
			s:      "this is a very long string",
			maxLen: 10,
			want:   "this is a ...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateString(tt.s, tt.maxLen); got != tt.want {
				t.Errorf("TruncateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{
			name:   "active status",
			status: "active",
		},
		{
			name:   "error status",
			status: "error",
		},
		{
			name:   "warning status",
			status: "warning",
		},
		{
			name:   "unknown status",
			status: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatStatus(tt.status)
			if result == "" {
				t.Error("FormatStatus should not return empty string")
			}
		})
	}
}
