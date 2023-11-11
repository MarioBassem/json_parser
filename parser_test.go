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
		"valid_with_space": {
			input:    "1234 ",
			expected: true,
			value:    1234,
		},
		"alphabetic": {
			input:    "abcd",
			expected: false,
			value:    0,
		},
		"alpha_numeric": {
			input:    "12ab34",
			expected: true,
			value:    12,
		},
		"plus_zero": {
			input:    "+0",
			expected: false,
			value:    0,
		},
		"negative_zero": {
			input:    "-0 ",
			expected: true,
			value:    0,
		},
		"multiple_decimal_points": {
			input:    "1.2.3",
			expected: false,
			value:    0,
		},
		"valid_int": {
			input:    "1234",
			expected: true,
			value:    1234,
		},
		"valid_float": {
			input:    "1.1234e5",
			expected: true,
			value:    1.1234e5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := parser{b: []byte(tc.input)}
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
			p := parser{b: []byte(tc.input)}
			str, err := p.getString()
			if tc.expected {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.value, str)
		})
	}
}

func TestGetValue(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected bool
		value    interface{}
	}{
		"empty": {
			input:    "",
			expected: false,
			value:    nil,
		},
		"not_ended_string": {
			input:    `\"abc`,
			expected: false,
			value:    nil,
		},
		"not_ended_array": {
			input:    `[1234,562`,
			expected: false,
			value:    nil,
		},
		"true_value": {
			input:    "true ",
			expected: true,
			value:    true,
		},
		"false_value": {
			input:    "false ",
			expected: true,
			value:    false,
		},
		"null_value": {
			input:    "null ",
			expected: true,
			value:    nil,
		},
		"unqoutoed_string": {
			input:    "hamada",
			expected: false,
			value:    nil,
		},
		"valid_string": {
			input:    `"hello world"`,
			expected: true,
			value:    "hello world",
		},
		"valid_array": {
			input:    `[1,2,3]`,
			expected: true,
			value:    []interface{}{float64(1), float64(2), float64(3)},
		},
		"valid_number": {
			input:    "-0.234e5 ",
			expected: true,
			value:    -0.234e5,
		},
		"valid_object": {
			input:    `{"key1": "value1"}`,
			expected: true,
			value:    map[string]interface{}{"key1": "value1"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := parser{b: []byte(tc.input)}
			val, err := p.getValue()
			if tc.expected {
				assert.NoError(t, err)
				assert.Equal(t, tc.value, val)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestGetArray(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected bool
		value    []interface{}
	}{
		"empty": {
			input:    "",
			expected: false,
			value:    nil,
		},
		"started_not_ended": {
			input:    "[1,2,3",
			expected: false,
			value:    nil,
		},
		"not_started_ended": {
			input:    "1,2,3]",
			expected: false,
			value:    nil,
		},
		"missing_comma": {
			input:    "[1,2{\"hello\":\"world\"}]",
			expected: false,
			value:    nil,
		},
		"valid_number_array": {
			input:    "[1,2,3]",
			expected: true,
			value:    []interface{}{float64(1), float64(2), float64(3)},
		},
		"valid_string_array": {
			input:    `["hello", "world"]`,
			expected: true,
			value:    []interface{}{"hello", "world"},
		},
		"valid_object_array": {
			input:    `[{\t"key1": "val1"},    {   "key2": 1234}]`,
			expected: true,
			value:    []interface{}{map[string]interface{}{"key1": "val1"}, map[string]interface{}{"key2": float64(1234)}},
		},
		"valid_mixed_array": {
			input:    `[1234, "string value", true, null]`,
			expected: true,
			value:    []interface{}{float64(1234), "string value", true, nil},
		},
		"valid_nested_array": {
			input:    `[1,2,[1,2,3]]`,
			expected: true,
			value:    []interface{}{float64(1), float64(2), []interface{}{float64(1), float64(2), float64(3)}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := parser{b: []byte(tc.input)}
			val, err := p.getArray()
			if tc.expected {
				assert.NoError(t, err)
				assert.Equal(t, tc.value, val)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
