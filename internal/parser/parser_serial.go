package parser

import (
	"fmt"
	"github.com/tablerst/jsonhelper/internal/lexer"
	"log"
	"strconv"
)

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peerToken
	p.peerToken = p.l.NextToken()
}

func (p *Parser) Parser() *ParseResult {
	result := &ParseResult{}

	switch p.curToken.Type {
	case lexer.TokenCurlyBraceOpen:
		result.Root = p.parseObject()
	case lexer.TokenSquareBracketOpen:
		result.Root = p.parseArray()
	default:
		result.Errors = append(result.Errors, fmt.Sprintf("Unexpected token: %s", p.curToken.Literal))
	}

	return result
}

func (p *Parser) parseObject() *ObjectNode {
	log.Println("Parsing object")
	node := &ObjectNode{
		Pairs: make(map[string]Node),
	}

	if p.curToken.Type != lexer.TokenCurlyBraceOpen {
		log.Printf("Expected TokenCurlyBraceOpen, got %s", p.curToken.Literal)
		return nil
	}

	p.nextToken()

	//var pendingComments []string

	for p.curToken.Type != lexer.TokenCurlyBraceClose && p.curToken.Type != lexer.TokenEOF {
		// Skip comments
		if p.curToken.Type == lexer.TokenComment {
			//pendingComments = append(pendingComments, p.curToken.Literal)
			p.nextToken()
			continue
		}

		key := p.curToken.Literal
		log.Printf("Parsing key: %s", key)

		p.nextToken()
		if p.curToken.Type != lexer.TokenColon {
			log.Printf("Expected TokenColon after key, got %s", p.curToken.Literal)
			return nil
		}

		p.nextToken()
		value := p.parseValue()
		if value == nil {
			log.Printf("Failed to parse value for key: %s", key)
			return nil
		}

		node.Pairs[key] = value

		p.nextToken()
		if p.curToken.Type == lexer.TokenComma {
			p.nextToken()
		}
	}

	if p.curToken.Type != lexer.TokenCurlyBraceClose {
		log.Printf("Expected TokenCurlyBraceClose, got %s", p.curToken.Literal)
		return nil
	}

	log.Println("Successfully parsed object")
	return node
}

func (p *Parser) parseArray() *ArrayNode {
	node := &ArrayNode{
		Elements: []Node{},
	}

	if p.curToken.Type != lexer.TokenSquareBracketOpen {
		return nil
	}

	p.nextToken()

	for p.curToken.Type != lexer.TokenSquareBracketClose && p.curToken.Type != lexer.TokenEOF {
		element := p.parseValue()
		if element == nil {
			return nil
		}
		node.Elements = append(node.Elements, element)

		p.nextToken()
		if p.curToken.Type == lexer.TokenComma {
			p.nextToken()
		}
	}

	if p.curToken.Type != lexer.TokenSquareBracketClose {
		return nil
	}

	return node
}

func (p *Parser) parseValue() Node {
	switch p.curToken.Type {
	case lexer.TokenString:
		return &StringNode{Value: p.curToken.Literal}
	case lexer.TokenNumber:
		// Convert string to float64
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil
		}
		return &NumberNode{Value: value}
	case lexer.TokenBoolean:
		value := p.curToken.Literal == "true"
		return &BoolNode{Value: value}
	case lexer.TokenNull:
		return &NullNode{}
	case lexer.TokenInfinity:
		return &InfinityNode{Positive: true}
	case lexer.TokenNaN:
		return &NaNNode{}
	case lexer.TokenCurlyBraceOpen:
		return p.parseObject()
	case lexer.TokenSquareBracketOpen:
		return p.parseArray()
	default:
		return nil
	}
}
