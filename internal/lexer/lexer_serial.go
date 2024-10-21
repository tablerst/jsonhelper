package lexer

import (
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

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.currentChar {
	case '-':
		if isDigit(l.peekChar()) {
			tok.Type = TokenNumber
			tok.Literal = l.readNumberWithSign()
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
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
		tok.Literal = l.readString(l.currentChar)
		tok.Type = TokenString
		tok.Line = l.line
		tok.Column = l.column
		l.readChar()
		return tok
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok.Literal = l.readSingleLineComment()
			tok.Type = TokenComment
			tok.Line = l.line
			tok.Column = l.column
			l.readChar()
			return tok
		} else if l.peekChar() == '*' {
			l.readChar()
			tok.Literal = l.readMultiLineComment()
			tok.Type = TokenComment
			tok.Line = l.line
			tok.Column = l.column
			l.readChar()
			return tok
		}
	case 0:
		tok.Literal = ""
		tok.Type = TokenEOF
	default:
		if isDigit(l.currentChar) {
			tok.Literal = l.readNumber()
			tok.Type = TokenNumber
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else if unicode.IsLetter(rune(l.currentChar)) || l.currentChar == '_' {
			literal := l.readIdentifier()
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
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			tok = Token{Type: TokenEOF, Literal: string(l.currentChar), Line: l.line, Column: l.column}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString(quote byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.currentChar == quote || l.currentChar == 0 {
			break
		}
		// Skip escaped characters
		if l.currentChar == '\\' && l.peekChar() == quote {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.currentChar) {
		l.readChar()
	}
	if l.currentChar == '.' {
		l.readChar()
		for isDigit(l.currentChar) {
			l.readChar()
		}
	}
	if l.currentChar == 'e' || l.currentChar == 'E' {
		l.readChar()
		if l.currentChar == '+' || l.currentChar == '-' {
			l.readChar()
		}
		for isDigit(l.currentChar) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumberWithSign() string {
	position := l.position
	if l.currentChar == '-' {
		l.readChar()
	}
	for isDigit(l.currentChar) {
		l.readChar()
	}
	if l.currentChar == '.' {
		l.readChar()
		for isDigit(l.currentChar) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
func (l *Lexer) readSingleLineComment() string {
	position := l.position + 1
	for l.currentChar != '\n' && l.currentChar != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readMultiLineComment() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.currentChar == '*' && l.peekChar() == '/' {
			break
		}
		if l.currentChar == 0 { // Reached end of input without finding end of comment
			break
		}
	}
	literal := l.input[position:l.position]
	l.readChar() // Read the closing '/'
	//fmt.Printf("Parsed multi-line comment: [%s]\n", literal)
	return literal
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for unicode.IsLetter(rune(l.currentChar)) || isDigit(l.currentChar) || l.currentChar == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}
