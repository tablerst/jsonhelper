package lexer

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
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	currentChar  byte
	line         int
	column       int
}
