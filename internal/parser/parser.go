package parser

import (
	"github.com/tablerst/jsonhelper/internal/lexer"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peerToken lexer.Token
}

type ParseResult struct {
	Root   Node
	Errors []string
}
