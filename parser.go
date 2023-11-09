package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrEndOfText           = errors.New("EOT")
	ErrUnExpectedEndOfText = errors.New("unexpected end of text")
	ErrInvalidCharacter    = errors.New("invalid Character")

	trueVal  = []byte("true")
	falseVal = []byte("false")
	nullVal  = []byte("null")
)

type Parser struct {
	b []byte
}

func Parse(b []byte) (map[string]interface{}, error) {
	p := Parser{b: b}
	return p.getObject()
}

func (p *Parser) getObject() (map[string]interface{}, error) {
	// get openning curly brace
	// whitespace
	// get key
	// whitespace
	// comma
	// value
	// closing curly brace
	obj := map[string]interface{}{}
	if err := p.skipByte('{'); err != nil {
		return nil, err
	}

	for {
		key, err := p.getString()
		if err != nil {
			return nil, err
		}

		if err := p.skipByte(':'); err != nil {
			return nil, err
		}

		val, err := p.getValue()
		if err != nil {
			return nil, err
		}

		obj[string(key)] = val

		err = p.skipByte(',')
		if errors.Is(err, ErrInvalidCharacter) {
			// next characters is not comma, we can break
			break
		}
		if err != nil {
			return nil, err
		}
	}

	if err := p.skipByte('}'); err != nil {
		return nil, err
	}

	return obj, nil
}

func (p *Parser) getString() (string, error) {
	p.skipWhiteSpace()

	if err := p.skipByte('"'); err != nil {
		return "", err
	}

	key := strings.Builder{}
	for {
		c, escaped, err := p.getChar()
		if err != nil {
			return "", err
		}

		if c == "\"" && !escaped {
			break
		}

		_, _ = key.WriteString(c)
	}

	return key.String(), nil
}

func (p *Parser) skipWhiteSpace() {
	p.b = bytes.TrimLeft(p.b, " \t\r\n")
}

func (p *Parser) skipByte(b byte) error {
	if len(p.b) == 0 {
		return fmt.Errorf("%w: looking for '%c'", ErrUnExpectedEndOfText, b)
	}

	if p.b[0] != b {
		return fmt.Errorf("%w: looking for '%c", ErrInvalidCharacter, b)
	}

	p.b = p.b[1:]
	return nil
}

func (p *Parser) getChar() (string, bool, error) {
	if len(p.b) == 0 {
		return "", false, ErrUnExpectedEndOfText
	}

	c := p.b[0]
	p.b = p.b[1:]
	if c != '\\' {
		return string(c), false, nil
	}

	escapedVal, err := p.getEscapedValue()
	if err != nil {
		return "", false, err
	}

	return escapedVal, true, nil
}

// getEscapedValue validates escape sequence, and returns the escaped value
func (p *Parser) getEscapedValue() (string, error) {
	if len(p.b) == 0 {
		return "", ErrUnExpectedEndOfText
	}

	b := p.b[0]
	p.b = p.b[1:]

	switch b {
	case '"':
		return "\"", nil
	case '\\':
		return "\\", nil
	case '/':
		return "/", nil
	case 'b':
		return "\b", nil
	case 'f':
		return "\f", nil
	case 'n':
		return "\n", nil
	case 'r':
		return "\r", nil
	case 't':
		return "\t", nil
	case 'u':
		hexCode, err := p.getHexCode()
		if err != nil {
			return "", err
		}

		return string(rune(hexCode)), nil
	default:
		return "", fmt.Errorf("invalid character '%c' in string escape code", b)
	}
}

func (p *Parser) getHexCode() (uint16, error) {
	if len(p.b) < 4 {
		return 0, ErrUnExpectedEndOfText
	}

	code := p.b[0:4]
	p.b = p.b[4:]

	decimalValue, err := strconv.ParseUint(string(code), 16, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid hex value: %w", err)
	}

	return uint16(decimalValue), nil
}

func (p *Parser) getValue() (interface{}, error) {
	p.skipWhiteSpace()

	if len(p.b) == 0 {
		return nil, ErrUnExpectedEndOfText
	}

	if p.b[0] == '"' {
		return p.getString()
	}

	if p.b[0] == '{' {
		return p.getObject()
	}

	if p.b[0] == '[' {
		return p.getArray()
	}

	if p.b[0] == '-' || unicode.IsDigit(rune(p.b[0])) {
		return p.getNumber()
	}

	if p.canSkipVal(trueVal) {
		p.b = p.b[len(trueVal):]
		return true, nil
	}

	if p.canSkipVal(falseVal) {
		p.b = p.b[len(falseVal):]
		return false, nil
	}

	if p.canSkipVal(nullVal) {
		p.b = p.b[len(nullVal):]
		return nil, nil
	}

	return nil, fmt.Errorf("%w: unexpected '%c'", ErrInvalidCharacter, p.b[0])
}

func (p *Parser) canSkipVal(val []byte) bool {
	if len(p.b) <= len(val) {
		return false
	}

	for idx := range val {
		if val[idx] != p.b[idx] {
			return false
		}
	}

	switch p.b[len(val)] {
	case '}', ']', ',', ' ', '\t', '\n', '\r':
		return true
	}

	return false
}

func (p *Parser) getArray() ([]interface{}, error) {
	arr := []interface{}{}
	if err := p.skipByte('['); err != nil {
		return nil, err
	}

	p.skipWhiteSpace()

	err := p.skipByte(']')
	if err == nil {
		return arr, nil
	}

	for {
		val, err := p.getValue()
		if err != nil {
			return nil, err
		}

		arr = append(arr, val)

		err = p.skipByte(',')
		if errors.Is(err, ErrInvalidCharacter) {
			// next characters is not comma, we can break
			break
		}
		if err != nil {
			return nil, err
		}
	}

	if err := p.skipByte(']'); err != nil {
		return nil, err
	}

	return arr, nil
}

func (p *Parser) getNumber() (float64, error) {
	numStr := strings.Builder{}
	for {
		c, _, err := p.getChar()
		if err != nil {
			return 0, err
		}

		_, _ = numStr.WriteString(c)

		if len(p.b) == 0 {
			return 0, ErrUnExpectedEndOfText
		}

		next := p.b[0]
		if unicode.IsDigit(rune(next)) || next == 'e' || next == 'E' || next == '.' || next == '-' || next == '+' {
			continue
		}

		break
	}

	num, err := strconv.ParseFloat(numStr.String(), 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}