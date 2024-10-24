package lexer

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenCurlyBraceOpen
	TokenCurlyBraceClose
	TokenSquareBracketOpen
	TokenSquareBracketClose
	TokenComma
	TokenColon
	TokenString
	TokenMinus
	TokenNumber
	TokenBoolean
	TokenNull
	TokenComment
	TokenInfinity
	TokenNaN
	TokenNewLine
	TokenWhitespace
	TokenError
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
	Err     error
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	currentChar  byte
	line         int
	column       int
}

// lexer.go

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "TokenEOF"
	case TokenCurlyBraceOpen:
		return "TokenCurlyBraceOpen"
	case TokenCurlyBraceClose:
		return "TokenCurlyBraceClose"
	case TokenSquareBracketOpen:
		return "TokenSquareBracketOpen"
	case TokenSquareBracketClose:
		return "TokenSquareBracketClose"
	case TokenComma:
		return "TokenComma"
	case TokenColon:
		return "TokenColon"
	case TokenString:
		return "TokenString"
	case TokenMinus:
		return "TokenMinus"
	case TokenNumber:
		return "TokenNumber"
	case TokenBoolean:
		return "TokenBoolean"
	case TokenNull:
		return "TokenNull"
	case TokenComment:
		return "TokenComment"
	case TokenInfinity:
		return "TokenInfinity"
	case TokenNaN:
		return "TokenNaN"
	case TokenNewLine:
		return "TokenNewLine"
	case TokenWhitespace:
		return "TokenWhitespace"
	case TokenError:
		return "TokenError"
	default:
		return fmt.Sprintf("UnknownToken(%d)", int(t))
	}
}
