package parser

import (
	"errors"
	"unicode"
)

const (
	doubleQuotes = iota
	openingCurlyBraces
	closingCurlyBraces
	openingBrackets
	closingBrackets
	colon
	comma
	escapedQuotes
	space
	backSlash
	forwardSlash
	// escapedSolidus
	// escapedReverseSolidus
	// backspace
	// formfeed
	// linefeed
	// carriageReturn
	// horizontalTab
	// hexEscapeSequence
	negativeSign
	positiveSign
	decimalPoint
	// eNotation
	digit
	other
)

var ErrEndOfText = errors.New("EOT")

type Tokenizer struct {
	b []byte
}

type Token struct {
	Kind  int
	Value byte
}

func (t *Tokenizer) Next() (Token, error) {
	if len(t.b) == 0 {
		return Token{}, ErrEndOfText
	}

	ch := t.b[0]
	t.b = t.b[1:]
	token := Token{
		Value: ch,
	}

	if unicode.IsDigit(rune(ch)) {
		token.Kind = digit
		return token, nil
	}

	switch ch {
	case '"':
		token.Kind = doubleQuotes
	case '{':
		token.Kind = openingCurlyBraces
	case '}':
		token.Kind = closingCurlyBraces
	case '[':
		token.Kind = openingBrackets
	case ']':
		token.Kind = closingBrackets
	case ':':
		token.Kind = colon
	case ',':
		token.Kind = comma
	// case '\"':
	// 	token.kind = escapedQuotes
	case ' ':
		token.Kind = space
	case '\\':
		token.Kind = backSlash
	case '/':
		token.Kind = forwardSlash
	// case '\b':
	// 	token.Kind = backspace
	// case '\f':
	// 	token.Kind = formfeed
	// case '\n':
	// 	token.Kind = linefeed
	// case '\r':
	// 	token.Kind = carriageReturn
	// case '\t':
	// 	token.Kind = horizontalTab
	// case '\u':
	// 	token.kind = hexEscapeSequence
	case '-':
		token.Kind = negativeSign
	case '+':
		token.Kind = positiveSign
	case '.':
		token.Kind = decimalPoint
	// case 'e', 'E':
	// 	token.Kind = eNotation
	default:
		token.Kind = other
	}

	return token, nil
}
