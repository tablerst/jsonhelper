package lexer

import (
	"fmt"
	"strconv"
	"unicode"
)

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
	if l.currentChar == '\n' {
		l.line++
		l.column = 0
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	// Handle whitespace
	if isWhitespace(l.currentChar) {
		startLine := l.line
		startColumn := l.column
		literal := l.readWhitespace()
		tok = Token{Type: TokenWhitespace, Literal: literal, Line: startLine, Column: startColumn}
		return tok
	}

	// Handle newlines
	if l.currentChar == '\n' {
		startLine := l.line
		startColumn := l.column
		literal := l.readNewlines()
		tok = Token{Type: TokenNewLine, Literal: literal, Line: startLine, Column: startColumn}
		return tok
	}

	switch l.currentChar {
	case '-':
		if l.peekString(len("Infinity")) == "Infinity" {
			l.readChar()                     // Consume '-'
			identifier := l.readIdentifier() // Read 'Infinity'
			literal := "-" + identifier
			tok.Type = TokenInfinity
			tok.Literal = literal
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else if isDigit(l.peekChar()) {
			l.readChar() // Consume '-'
			num, err := l.readNumber()
			tok.Line = l.line
			tok.Column = l.column
			if err != nil {
				tok.Type = TokenError
				tok.Literal = ""
				tok.Err = err
			} else {
				tok.Literal = "-" + num
				tok.Type = TokenNumber
			}
			return tok
		} else {
			// In JSON5, minus can be a standalone token
			tok = Token{Type: TokenMinus, Literal: "-", Line: l.line, Column: l.column}
		}
	case '{':
		tok = Token{Type: TokenCurlyBraceOpen, Literal: "{", Line: l.line, Column: l.column}
	case '}':
		tok = Token{Type: TokenCurlyBraceClose, Literal: "}", Line: l.line, Column: l.column}
	case '[':
		tok = Token{Type: TokenSquareBracketOpen, Literal: "[", Line: l.line, Column: l.column}
	case ']':
		tok = Token{Type: TokenSquareBracketClose, Literal: "]", Line: l.line, Column: l.column}
	case ',':
		tok = Token{Type: TokenComma, Literal: ",", Line: l.line, Column: l.column}
	case ':':
		tok = Token{Type: TokenColon, Literal: ":", Line: l.line, Column: l.column}
	case '"', '\'':
		str, err := l.readString(l.currentChar)
		tok.Line = l.line
		tok.Column = l.column
		if err != nil {
			tok.Type = TokenError
			tok.Literal = ""
			tok.Err = err
		} else {
			tok.Type = TokenString
			tok.Literal = str
		}
		return tok
	case '/':
		if l.peekChar() == '/' {
			startPos := l.position // Include the '//' in the literal
			tok.Line = l.line
			tok.Column = l.column
			l.readChar() // Consume '/'
			l.readChar() // Consume second '/'
			tok.Literal = l.readSingleLineComment(startPos)
			tok.Type = TokenComment
			return tok
		} else if l.peekChar() == '*' {
			startPos := l.position // Include the '/*' in the literal
			tok.Line = l.line
			tok.Column = l.column
			l.readChar() // Consume '/'
			l.readChar() // Consume '*'
			literal, err := l.readMultiLineComment(startPos)
			if err != nil {
				tok.Type = TokenError
				tok.Literal = ""
				tok.Err = err
			} else {
				tok.Literal = literal
				tok.Type = TokenComment
			}
			return tok
		} else {
			// Handle division operator or invalid token
			tok = Token{Type: TokenError, Literal: string(l.currentChar), Line: l.line, Column: l.column, Err: fmt.Errorf("unexpected character: %q", l.currentChar)}
		}
	case 0:
		tok.Literal = ""
		tok.Type = TokenEOF
	default:
		if isDigit(l.currentChar) {
			num, err := l.readNumber()
			tok.Line = l.line
			tok.Column = l.column
			if err != nil {
				tok.Type = TokenError
				tok.Literal = ""
				tok.Err = err
			} else {
				tok.Type = TokenNumber
				tok.Literal = num
			}
			return tok
		} else if unicode.IsLetter(rune(l.currentChar)) || l.currentChar == '_' {
			literal := l.readIdentifier()
			tok.Line = l.line
			tok.Column = l.column
			switch literal {
			case "true", "false":
				tok.Type = TokenBoolean
			case "null":
				tok.Type = TokenNull
			case "Infinity":
				tok.Type = TokenInfinity
			case "NaN":
				tok.Type = TokenNaN
			default:
				tok.Type = TokenString
			}
			tok.Literal = literal
			return tok
		} else {
			// Unknown character
			tok = Token{Type: TokenError, Literal: string(l.currentChar), Line: l.line, Column: l.column, Err: fmt.Errorf("unknown character: %q", l.currentChar)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString(quote byte) (string, error) {
	var result []rune
	l.readChar() // Move past the opening quote
	for l.currentChar != quote {
		if l.currentChar == 0 {
			// Reached end of input without closing quote
			return "", fmt.Errorf("unterminated string literal")
		}
		if l.currentChar == '\\' {
			l.readChar()
			switch l.currentChar {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '"', '\'', '\\', '/':
				result = append(result, rune(l.currentChar))
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'u':
				// Handle Unicode escape sequences
				unicodeSeq, err := l.readUnicodeSequence()
				if err != nil {
					return "", err
				}
				result = append(result, unicodeSeq)
			default:
				// Handle other escape sequences or errors
				return "", fmt.Errorf("invalid escape sequence: \\%c", l.currentChar)
			}
		} else {
			result = append(result, rune(l.currentChar))
		}
		l.readChar()
	}
	l.readChar() // Consume closing quote
	return string(result), nil
}

func (l *Lexer) readUnicodeSequence() (rune, error) {
	if l.peekChar() == '{' {
		// Handle \u{XXXXXX}
		l.readChar() // Consume '{'
		hexDigits := ""
		for {
			l.readChar()
			if l.currentChar == '}' {
				break
			}
			if !isHexDigit(l.currentChar) {
				// Handle error
				return 0, fmt.Errorf("invalid Unicode escape sequence")
			}
			hexDigits += string(l.currentChar)
		}
		// Convert hexDigits to rune
		if codePoint, err := strconv.ParseInt(hexDigits, 16, 32); err == nil {
			if codePoint > 0x10FFFF {
				return 0, fmt.Errorf("Unicode code point out of range")
			}
			return rune(codePoint), nil
		} else {
			return 0, err
		}

	} else {
		// Handle \uXXXX
		hexDigits := ""
		for i := 0; i < 4; i++ {
			l.readChar()
			if !isHexDigit(l.currentChar) {
				// Handle error
				return 0, fmt.Errorf("invalid Unicode escape sequence")
			}
			hexDigits += string(l.currentChar)
		}
		// Convert hexDigits to rune
		if codePoint, err := strconv.ParseInt(hexDigits, 16, 16); err == nil {
			return rune(codePoint), nil
		} else {
			return 0, err
		}
	}
}

func (l *Lexer) readNumber() (string, error) {
	position := l.position
	// Handle leading '0' for hex, octal, binary
	if l.currentChar == '0' {
		if l.peekChar() == 'x' || l.peekChar() == 'X' {
			l.readChar() // Consume '0'
			l.readChar() // Consume 'x' or 'X'

			if !isHexDigit(l.currentChar) {
				return "", fmt.Errorf("invalid hexadecimal number: expected hex digit after '0x'")
			}
			for {
				if isHexDigit(l.currentChar) {
					l.readChar()
				} else if unicode.IsLetter(rune(l.currentChar)) || isDigit(l.currentChar) {
					return "", fmt.Errorf("invalid digit '%c' in hexadecimal number", l.currentChar)
				} else {
					break
				}
			}
			return l.input[position:l.position], nil
		} else if l.peekChar() == 'b' || l.peekChar() == 'B' {
			l.readChar() // Consume '0'
			l.readChar() // Consume 'b' or 'B'

			if !isBinaryDigit(l.currentChar) {
				return "", fmt.Errorf("invalid binary number: expected binary digit after '0b'")
			}
			for {
				if isBinaryDigit(l.currentChar) {
					l.readChar()
				} else if isDigit(l.currentChar) || unicode.IsLetter(rune(l.currentChar)) {
					return "", fmt.Errorf("invalid digit '%c' in binary number", l.currentChar)
				} else {
					break
				}
			}
			return l.input[position:l.position], nil
		} else if l.peekChar() == 'o' || l.peekChar() == 'O' {
			l.readChar() // Consume '0'
			l.readChar() // Consume 'o' or 'O'

			if !isOctalDigit(l.currentChar) {
				return "", fmt.Errorf("invalid octal number: expected octal digit after '0o'")
			}
			for {
				if isOctalDigit(l.currentChar) {
					l.readChar()
				} else if isDigit(l.currentChar) || unicode.IsLetter(rune(l.currentChar)) {
					return "", fmt.Errorf("invalid digit '%c' in octal number", l.currentChar)
				} else {
					break
				}
			}
			return l.input[position:l.position], nil
		}
	}

	// Decimal number parsing
	if !isDigit(l.currentChar) && l.currentChar != '.' {
		return "", fmt.Errorf("invalid number")
	}

	for isDigit(l.currentChar) {
		l.readChar()
	}

	if l.currentChar == '.' {
		l.readChar()
		if !isDigit(l.currentChar) {
			return "", fmt.Errorf("invalid decimal number")
		}
		for isDigit(l.currentChar) {
			l.readChar()
		}
	}

	if l.currentChar == 'e' || l.currentChar == 'E' {
		l.readChar()
		if l.currentChar == '+' || l.currentChar == '-' {
			l.readChar()
		}
		if !isDigit(l.currentChar) {
			return "", fmt.Errorf("invalid exponent in number")
		}
		for isDigit(l.currentChar) {
			l.readChar()
		}
	}

	return l.input[position:l.position], nil
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for unicode.IsLetter(rune(l.currentChar)) || isDigit(l.currentChar) || l.currentChar == '_' || l.currentChar == '$' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readSingleLineComment(startPos int) string {
	for l.currentChar != '\n' && l.currentChar != 0 {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) readMultiLineComment(startPos int) (string, error) {
	for {
		if l.currentChar == 0 {
			// Reached end of input without closing comment
			return "", fmt.Errorf("unterminated multi-line comment")
		}
		if l.currentChar == '*' && l.peekChar() == '/' {
			l.readChar() // Consume '*'
			l.readChar() // Consume '/'
			break
		}
		l.readChar()
	}
	return l.input[startPos:l.position], nil
}

func (l *Lexer) peekString(n int) string {
	endPos := l.readPosition + n - 1
	if endPos > len(l.input) {
		endPos = len(l.input)
	}
	return l.input[l.readPosition:endPos]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

func isBinaryDigit(ch byte) bool {
	return ch == '0' || ch == '1'
}

func isOctalDigit(ch byte) bool {
	return ch >= '0' && ch <= '7'
}

func (l *Lexer) readWhitespace() string {
	position := l.position
	for isWhitespace(l.currentChar) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNewlines() string {
	position := l.position
	for l.currentChar == '\n' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}
