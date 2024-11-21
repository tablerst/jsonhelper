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
	LBRACE     // '{'
	RBRACE     // '}'
	LBRACKET   // '['
	RBRACKET   // ']'
	COLON      // ':'
	COMMA      // ','
	IDENTIFIER // 未加引号的键名
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
		l.readChar()
	case '}':
		tok = Token{Type: RBRACE, Value: string(l.ch)}
		l.readChar()
	case '[':
		tok = Token{Type: LBRACKET, Value: string(l.ch)}
		l.readChar()
	case ']':
		tok = Token{Type: RBRACKET, Value: string(l.ch)}
		l.readChar()
	case ':':
		tok = Token{Type: COLON, Value: string(l.ch)}
		l.readChar()
	case ',':
		tok = Token{Type: COMMA, Value: string(l.ch)}
		l.readChar()
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
		} else if isDigit(l.ch) || l.ch == '-' || l.ch == '+' || l.ch == '.' {
			number := l.readNumber()
			tok = Token{Type: NUMBER, Value: number}
		} else {
			tok = Token{Type: ERROR, Value: string(l.ch)}
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
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
	l.readChar()
	var sb strings.Builder
	for {
		if l.ch == quote {
			break
		} else if l.ch == '\\' {
			l.readChar()
			if l.ch == '\n' || l.ch == '\r' {
				l.skipLineBreak()
			} else {
				sb.WriteRune('\\')
				sb.WriteRune(l.ch)
				l.readChar()
			}
		} else if l.ch == 0 {
			// 未终止的字符串
			break
		} else {
			sb.WriteRune(l.ch)
			l.readChar()
		}
	}
	l.readChar() // 读取结束引号
	return sb.String()
}

func (l *Lexer) skipLineBreak() {
	if l.ch == '\r' {
		l.readChar()
		if l.ch == '\n' {
			l.readChar()
		}
	} else if l.ch == '\n' {
		l.readChar()
	}
	// 跳过换行符后的空白字符
	l.skipWhitespace()
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_' || ch == '$'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position

	if l.ch == '-' || l.ch == '+' {
		l.readChar()
	}

	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		// 处理十六进制数字
		l.readChar()
		l.readChar()
		for isHexDigit(l.ch) {
			l.readChar()
		}
	} else {
		decimalPointFound := false

		if l.ch == '.' {
			decimalPointFound = true
			l.readChar()
		}

		for isDigit(l.ch) {
			l.readChar()
		}

		if l.ch == '.' && !decimalPointFound {
			decimalPointFound = true
			l.readChar()
			for isDigit(l.ch) {
				l.readChar()
			}
		}

		if l.ch == 'e' || l.ch == 'E' {
			l.readChar()
			if l.ch == '+' || l.ch == '-' {
				l.readChar()
			}
			for isDigit(l.ch) {
				l.readChar()
			}
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
		return r
	}
}

func isHexDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
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
		return IDENTIFIER
	}
}
