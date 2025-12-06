package unixtime2rfc

import (
	"errors"
	"testing"
	"time"
)

// Helper to create a pointer to an int64.
func int64p(i int64) *int64 {
	return &i
}

func TestProcessTimeInput(t *testing.T) {
	testCases := []struct {
		name        string
		input       UnixTimeInput
		expected    time.Time
		expectErr   bool
		expectedErr error
	}{
		{
			name:      "UnixTime only - epoch",
			input:     UnixTimeInput{UnixTime: int64p(0)},
			expected:  time.Unix(0, 0).UTC(),
			expectErr: false,
		},
		{
			name:      "UnixTime only - specific",
			input:     UnixTimeInput{UnixTime: int64p(1678886400)},
			expected:  time.Unix(1678886400, 0).UTC(),
			expectErr: false,
		},
		{
			name:      "UnixTimeMs only - epoch",
			input:     UnixTimeInput{UnixTimeMs: int64p(0)},
			expected:  time.UnixMilli(0).UTC(),
			expectErr: false,
		},
		{
			name:      "UnixTimeMs only - specific",
			input:     UnixTimeInput{UnixTimeMs: int64p(1678886400123)},
			expected:  time.UnixMilli(1678886400123).UTC(),
			expectErr: false,
		},
		{
			name:      "UnixTimeUs only - epoch",
			input:     UnixTimeInput{UnixTimeUs: int64p(0)},
			expected:  time.UnixMicro(0).UTC(),
			expectErr: false,
		},
		{
			name:      "UnixTimeUs only - specific",
			input:     UnixTimeInput{UnixTimeUs: int64p(1678886400000456)},
			expected:  time.UnixMicro(1678886400000456).UTC(),
			expectErr: false,
		},
		{
			name:      "Precedence: UnixTime over Ms and Us",
			input:     UnixTimeInput{UnixTime: int64p(1), UnixTimeMs: int64p(2), UnixTimeUs: int64p(3)},
			expected:  time.Unix(1, 0).UTC(),
			expectErr: false,
		},
		{
			name:      "Precedence: Ms over Us",
			input:     UnixTimeInput{UnixTimeMs: int64p(2), UnixTimeUs: int64p(3)},
			expected:  time.UnixMilli(2).UTC(),
			expectErr: false,
		},
		{
			name:        "No timestamp provided",
			input:       UnixTimeInput{Layout: "RFC3339"},
			expected:    time.Time{},
			expectErr:   true,
			expectedErr: ErrNoTimestampProvided,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ProcessTimeInput(tc.input)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("Expected error '%v', but got '%v'", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Did not expect an error, but got: %v", err)
			}

			if !result.Equal(tc.expected) {
				t.Errorf("Expected time '%s', but got '%s'", tc.expected, result)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	// A non-zero timestamp to make the formats distinguishable
	const testTime int64 = 1678886400 // 2023-03-15T12:00:00Z
	testTimeObj := time.Unix(testTime, 0).UTC()

	testCases := []struct {
		name        string
		inputTime   time.Time
		layout      string
		expected    string
		expectErr   bool
		expectedErr error
	}{
		{
			name:      "Epoch time with default layout",
			inputTime: time.Unix(0, 0).UTC(),
			layout:    "",
			expected:  "1970-01-01T00:00:00Z",
			expectErr: false,
		},
		{
			name:      "Specific time with DateOnly layout",
			inputTime: testTimeObj,
			layout:    "DateOnly",
			expected:  "2023-03-15",
			expectErr: false,
		},
		{
			name:        "Invalid layout returns error",
			inputTime:   time.Unix(0, 0).UTC(),
			layout:      "NoSuchLayout",
			expected:    "",
			expectErr:   true,
			expectedErr: ErrInvalidLayout,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FormatTime(tc.inputTime, tc.layout)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected an error for layout '%s', but got none", tc.layout)
				}
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("Expected error '%v', but got '%v'", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Did not expect an error for layout '%s', but got: %v", tc.layout, err)
			}

			if result != tc.expected {
				t.Errorf("Expected '%s', but got '%s'", tc.expected, result)
			}
		})
	}
}
