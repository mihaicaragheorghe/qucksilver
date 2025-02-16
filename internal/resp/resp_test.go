package resp_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/mihaicaragheorghe/qucksilver/internal/resp"
)

func TestParseRESP(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "Simple String",
			input:    "+OK\r\n",
			expected: "OK",
			wantErr:  false,
		},
		{
			name:     "Error",
			input:    "-Error message\r\n",
			expected: "Error message",
			wantErr:  false,
		},
		{
			name:     "Integer",
			input:    ":1000\r\n",
			expected: int64(1000),
			wantErr:  false,
		},
		{
			name:     "Bulk String",
			input:    "$6\r\nfoobar\r\n",
			expected: "foobar",
			wantErr:  false,
		},
		{
			name:     "Array",
			input:    "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			expected: []interface{}{"foo", "bar"},
			wantErr:  false,
		},
		{
			name:     "Null Bulk String",
			input:    "$-1\r\n",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "Null Array",
			input:    "*-1\r\n",
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewReader([]byte(tc.input)))
			result, err := resp.ParseRESP(reader)

			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}

			if err == nil && !compare(result, tc.expected) {
				t.Errorf("expected: %#v, got: %#v", tc.expected, result)
			}
		})
	}
}

func compare(actual, expected interface{}) bool {
	if actual == nil && expected == nil {
		return true
	}
	// Handle typed nil slices vs untyped nil
	if actualSlice, ok := actual.([]interface{}); ok && expected == nil {
		return actualSlice == nil
	}

	switch exp := expected.(type) {
	case []interface{}:
		act, ok := actual.([]interface{})
		if !ok || len(act) != len(exp) {
			return false
		}
		for i := range act {
			if !compare(act[i], exp[i]) {
				return false
			}
		}
		return true
	default:
		return actual == expected
	}
}
