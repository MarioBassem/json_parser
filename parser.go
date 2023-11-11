package parser

import (
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

type parser struct {
	b []byte
}

func Parse(b []byte) (map[string]interface{}, error) {
	p := parser{b: b}
	return p.getObject()
}

func (p *parser) getObject() (map[string]interface{}, error) {
	// get openning curly brace
	// whitespace
	// get key
	// whitespace
	// comma
	// value
	// closing curly brace
	p.skipWhiteSpace()
	defer p.skipWhiteSpace()

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

func (p *parser) getString() (string, error) {
	p.skipWhiteSpace()
	defer p.skipWhiteSpace()

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

func (p *parser) skipWhiteSpace() {
	cut := 0
	for i := 0; i < len(p.b); i++ {
		c1 := p.b[i]
		if c1 == ' ' || c1 == '\t' || c1 == '\n' || c1 == '\r' {
			cut++
			continue
		}

		if len(p.b) == 1 {
			break
		}

		c2 := p.b[i+1]
		if c1 == '\\' && (c2 == 't' || c2 == 'n' || c2 == 'r') {
			cut += 2
			i++
			continue
		}

		break
	}

	p.b = p.b[cut:]
}

func (p *parser) skipByte(b byte) error {
	if len(p.b) == 0 {
		return fmt.Errorf("%w: looking for '%c'", ErrUnExpectedEndOfText, b)
	}

	if p.b[0] != b {
		return fmt.Errorf("%w: looking for '%c'", ErrInvalidCharacter, b)
	}

	p.b = p.b[1:]
	return nil
}

func (p *parser) getChar() (string, bool, error) {
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
func (p *parser) getEscapedValue() (string, error) {
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

func (p *parser) getHexCode() (uint16, error) {
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

func (p *parser) getValue() (interface{}, error) {
	p.skipWhiteSpace()
	defer p.skipWhiteSpace()

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

func (p *parser) canSkipVal(val []byte) bool {
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

func (p *parser) getArray() ([]interface{}, error) {
	p.skipWhiteSpace()
	defer p.skipWhiteSpace()

	arr := []interface{}{}
	if err := p.skipByte('['); err != nil {
		return nil, err
	}

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

func (p *parser) getNumber() (float64, error) {
	p.skipWhiteSpace()
	defer p.skipWhiteSpace()

	numStr := strings.Builder{}
	for {
		c, _, err := p.getChar()
		if err != nil {
			return 0, err
		}

		_, _ = numStr.WriteString(c)

		if len(p.b) == 0 {
			break
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
