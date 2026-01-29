package dateparse

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	// Reference time: Monday, January 15, 2024, 10:00 AM
	ref := time.Date(2024, 1, 15, 10, 0, 0, 0, time.Local)

	tests := []struct {
		name     string
		input    string
		ref      time.Time
		want     time.Time
		wantErr  bool
		errMatch string
	}{
		{
			name:  "today",
			input: "today",
			ref:   ref,
			want:  time.Date(2024, 1, 15, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "tomorrow",
			input: "tomorrow",
			ref:   ref,
			want:  time.Date(2024, 1, 16, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "yesterday",
			input: "yesterday",
			ref:   ref,
			want:  time.Date(2024, 1, 14, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "next week",
			input: "next week",
			ref:   ref,
			want:  time.Date(2024, 1, 22, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "next monday",
			input: "next monday",
			ref:   ref,
			want:  time.Date(2024, 1, 22, 17, 0, 0, 0, time.Local), // Next Monday from Jan 15 (Mon)
		},
		{
			name:  "next tuesday",
			input: "next tuesday",
			ref:   ref,
			want:  time.Date(2024, 1, 16, 17, 0, 0, 0, time.Local), // Next Tuesday from Jan 15 (Mon)
		},
		{
			name:  "next wednesday",
			input: "next wednesday",
			ref:   ref,
			want:  time.Date(2024, 1, 17, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "next thursday",
			input: "next thursday",
			ref:   ref,
			want:  time.Date(2024, 1, 18, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "next friday",
			input: "next friday",
			ref:   ref,
			want:  time.Date(2024, 1, 19, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "next saturday",
			input: "next saturday",
			ref:   ref,
			want:  time.Date(2024, 1, 20, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "next sunday",
			input: "next sunday",
			ref:   ref,
			want:  time.Date(2024, 1, 21, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "in 3 days",
			input: "in 3 days",
			ref:   ref,
			want:  time.Date(2024, 1, 18, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "in 1 day",
			input: "in 1 day",
			ref:   ref,
			want:  time.Date(2024, 1, 16, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "in 2 weeks",
			input: "in 2 weeks",
			ref:   ref,
			want:  time.Date(2024, 1, 29, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "in 1 week",
			input: "in 1 week",
			ref:   ref,
			want:  time.Date(2024, 1, 22, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "ISO format",
			input: "2024-03-20",
			ref:   ref,
			want:  time.Date(2024, 3, 20, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "short month day",
			input: "Jan 15",
			ref:   ref,
			want:  time.Date(2024, 1, 15, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "long month day",
			input: "January 15",
			ref:   ref,
			want:  time.Date(2024, 1, 15, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "short month day year",
			input: "Jan 15 2024",
			ref:   ref,
			want:  time.Date(2024, 1, 15, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "long month day year",
			input: "January 15 2024",
			ref:   ref,
			want:  time.Date(2024, 1, 15, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "march abbreviated",
			input: "Mar 5",
			ref:   ref,
			want:  time.Date(2024, 3, 5, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "december full name with year",
			input: "December 31 2025",
			ref:   ref,
			want:  time.Date(2025, 12, 31, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "case insensitive TODAY",
			input: "TODAY",
			ref:   ref,
			want:  time.Date(2024, 1, 15, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "case insensitive Next Monday",
			input: "Next Monday",
			ref:   ref,
			want:  time.Date(2024, 1, 22, 17, 0, 0, 0, time.Local),
		},
		{
			name:  "case insensitive IN 3 DAYS",
			input: "IN 3 DAYS",
			ref:   ref,
			want:  time.Date(2024, 1, 18, 17, 0, 0, 0, time.Local),
		},
		{
			name:     "empty string",
			input:    "",
			ref:      ref,
			wantErr:  true,
			errMatch: "empty",
		},
		{
			name:     "invalid format",
			input:    "invalid",
			ref:      ref,
			wantErr:  true,
			errMatch: "unrecognized",
		},
		{
			name:     "invalid number format",
			input:    "in abc days",
			ref:      ref,
			wantErr:  true,
			errMatch: "unrecognized",
		},
		{
			name:     "partial match",
			input:    "nextt monday",
			ref:      ref,
			wantErr:  true,
			errMatch: "unrecognized",
		},
		{
			name:     "invalid day in month-day format",
			input:    "Jan abc",
			ref:      ref,
			wantErr:  true,
			errMatch: "unrecognized",
		},
		{
			name:     "invalid year in month-day-year format",
			input:    "Jan 15 xyz",
			ref:      ref,
			wantErr:  true,
			errMatch: "unrecognized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseWithReference(tt.input, tt.ref)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseWithReference() expected error containing %q, got nil", tt.errMatch)
					return
				}
				if tt.errMatch != "" && !contains(err.Error(), tt.errMatch) {
					t.Errorf("ParseWithReference() error = %v, want error containing %q", err, tt.errMatch)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseWithReference() unexpected error = %v", err)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParseWithReference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse_UsesCurrentTime(t *testing.T) {
	// Parse without reference should use current time
	got, err := Parse("today")
	if err != nil {
		t.Fatalf("Parse() unexpected error = %v", err)
	}

	now := time.Now()
	expected := time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, time.Local)

	if !got.Equal(expected) {
		t.Errorf("Parse(\"today\") = %v, want %v", got, expected)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
