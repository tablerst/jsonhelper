package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	EOF TokenType = iota
	ERROR
	STRING
	NUMBER
	BOOLEAN
	NULL
	NAN
	INFINITY
	HEX
	LBRACE   // '{'
	RBRACE   // '}'
	LBRACKET // '['
	RBRACKET // ']'
	COLON    // ':'
	COMMA    // ','
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input        string
	position     int  // 当前输入的位置
	readPosition int  // 下一次读取的位置
	ch           rune // 当前字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		r, width := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += width
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '{':
		tok = Token{Type: LBRACE, Value: string(l.ch)}
	case '}':
		tok = Token{Type: RBRACE, Value: string(l.ch)}
	case '[':
		tok = Token{Type: LBRACKET, Value: string(l.ch)}
	case ']':
		tok = Token{Type: RBRACKET, Value: string(l.ch)}
	case ':':
		tok = Token{Type: COLON, Value: string(l.ch)}
	case ',':
		tok = Token{Type: COMMA, Value: string(l.ch)}
	case '"', '\'':
		tok = Token{Type: STRING, Value: l.readString(l.ch)}
	case 0:
		tok = Token{Type: EOF, Value: ""}
	case '/', '#':
		l.readComment()
		return l.NextToken()
	default:
		if isLetter(l.ch) {
			identifier := l.readIdentifier()
			tokType := lookupIdent(identifier)
			tok = Token{Type: tokType, Value: identifier}
			return tok
		} else if isDigit(l.ch) || l.ch == '-' || l.ch == '+' || l.ch == '.' {
			number := l.readNumber()
			tok = Token{Type: NUMBER, Value: number}
			return tok
		} else {
			tok = Token{Type: ERROR, Value: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readComment() {
	if l.ch == '/' {
		l.readChar()
		if l.ch == '/' {
			// 单行注释
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		} else if l.ch == '*' {
			// 多行注释
			l.readChar()
			for {
				if l.ch == '*' {
					l.readChar()
					if l.ch == '/' {
						l.readChar()
						break
					}
				} else if l.ch == 0 {
					break
				} else {
					l.readChar()
				}
			}
		} else {
			// 其他情况
		}
	} else if l.ch == '#' {
		// 单行注释
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}
}

func (l *Lexer) readString(quote rune) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == quote {
			break
		} else if l.ch == '\\' {
			l.readChar()
		} else if l.ch == 0 {
			break
		}
	}
	result := l.input[position:l.position]
	l.readChar()
	return result
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' || l.ch == 'e' || l.ch == 'E' || l.ch == '-' || l.ch == '+' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func lookupIdent(ident string) TokenType {
	switch strings.ToLower(ident) {
	case "true", "false":
		return BOOLEAN
	case "null":
		return NULL
	case "nan":
		return NUMBER
	case "infinity":
		return NUMBER
	default:
		return STRING
	}
}
