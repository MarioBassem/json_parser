package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestParse(t *testing.T) {

// }

func TestGetNumber(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected bool
		value    float64
	}{
		"empty": {
			input:    "",
			expected: false,
			value:    0,
		},
		"not_ended": {
			input:    "1234",
			expected: false,
			value:    0,
		},
		"alphabetic": {
			input:    "abcd ",
			expected: false,
			value:    0,
		},
		"alpha_numeric": {
			input:    "12ab34 ",
			expected: true,
			value:    12,
		},
		"plus_zero": {
			input:    "+0 ",
			expected: false,
			value:    0,
		},
		"negative_zero": {
			input:    "-0 ",
			expected: true,
			value:    0,
		},
		"multiple_decimal_points": {
			input:    "1.2.3 ",
			expected: false,
			value:    0,
		},
		"valid_int": {
			input:    "1234 ",
			expected: true,
			value:    1234,
		},
		"valid_float": {
			input:    "1.1234e5 ",
			expected: true,
			value:    1.1234e5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Parser{b: []byte(tc.input)}
			num, err := p.getNumber()
			if tc.expected {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.value, num)
		})
	}
}

func TestGetString(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected bool
		value    string
	}{
		"empty": {
			input:    ``,
			expected: false,
			value:    "",
		},
		"not_ended": {
			input:    `"hi`,
			expected: false,
			value:    "",
		},
		"not_started": {
			input:    `hii"`,
			expected: false,
			value:    "",
		},
		"with_valid_escape_sequences": {
			input:    `"hi\n\tmy \"friend. enough\u1234with the escape_sequences"`,
			expected: true,
			value:    "hi\n\tmy \"friend. enough\u1234with the escape_sequences",
		},
		"with_invalid_escape_sequences": {
			input:    `"hello\qworld"`,
			expected: false,
		},
		"valid": {
			input:    `"hello world"`,
			expected: true,
			value:    "hello world",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Parser{b: []byte(tc.input)}
			str, err := p.getString()
			if tc.expected {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.value, str)
		})
	}
}
