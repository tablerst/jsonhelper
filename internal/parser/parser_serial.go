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

	var leadingComments []*CommentNode

	// Collect leading comments
	for p.curToken.Type == lexer.TokenComment {
		commentNode := &CommentNode{
			JSONNode: &JSONNode{},
			Text:     p.curToken.Literal,
		}
		leadingComments = append(leadingComments, commentNode)
		p.nextToken()
	}

	var root Node

	switch p.curToken.Type {
	case lexer.TokenCurlyBraceOpen:
		root = p.parseObject()
	case lexer.TokenSquareBracketOpen:
		root = p.parseArray()
	default:
		result.Errors = append(result.Errors, fmt.Sprintf("Unexpected token: %s", p.curToken.Literal))
		return result
	}

	if root != nil {
		root.AddLeadingComments(leadingComments)
		result.Root = root
	} else {
		result.Errors = append(result.Errors, "Failed to parse root element")
	}

	return result
}

func (p *Parser) parseObject() *ObjectNode {
	log.Println("Parsing object")
	node := &ObjectNode{
		JSONNode: &JSONNode{},
		Pairs:    []KeyValuePair{},
	}

	if p.curToken.Type != lexer.TokenCurlyBraceOpen {
		log.Printf("Expected TokenCurlyBraceOpen, got %s", p.curToken.Literal)
		return nil
	}

	p.nextToken()

	for p.curToken.Type != lexer.TokenCurlyBraceClose && p.curToken.Type != lexer.TokenEOF {
		for p.curToken.Type == lexer.TokenNewLine {
			p.nextToken()
		}
		// 收集前置注释
		var leadingComments []*CommentNode
		for p.curToken.Type == lexer.TokenComment {
			commentNode := &CommentNode{
				JSONNode: &JSONNode{},
				Text:     p.curToken.Literal,
			}
			leadingComments = append(leadingComments, commentNode)
			p.nextToken()

			for p.curToken.Type == lexer.TokenNewLine {
				p.nextToken()
			}
		}

		// 解析键
		if p.curToken.Type != lexer.TokenString {
			log.Printf("Expected TokenString as key, got %s", p.curToken.Literal)
			return nil
		}
		key := p.curToken.Literal
		log.Printf("Parsing key: %s", key)
		p.nextToken()

		// 收集键后的注释
		var keyTrailingComments []*CommentNode
		for p.curToken.Type == lexer.TokenComment {
			commentNode := &CommentNode{
				JSONNode: &JSONNode{},
				Text:     p.curToken.Literal,
			}
			keyTrailingComments = append(keyTrailingComments, commentNode)
			p.nextToken()
		}

		if p.curToken.Type != lexer.TokenColon {
			log.Printf("Expected TokenColon after key, got %s", p.curToken.Literal)
			return nil
		}
		p.nextToken()

		// 收集值的前置注释
		var valueLeadingComments []*CommentNode
		for p.curToken.Type == lexer.TokenComment {
			commentNode := &CommentNode{
				JSONNode: &JSONNode{},
				Text:     p.curToken.Literal,
			}
			valueLeadingComments = append(valueLeadingComments, commentNode)
			p.nextToken()
		}

		// 解析值
		value := p.parseValue()
		if value == nil {
			log.Printf("Failed to parse value for key: %s", key)
			return nil
		}
		value.AddLeadingComments(valueLeadingComments)

		// 收集值后的注释
		var valueTrailingComments []*CommentNode
		for p.curToken.Type == lexer.TokenComment {
			commentNode := &CommentNode{
				JSONNode: &JSONNode{},
				Text:     p.curToken.Literal,
			}
			valueTrailingComments = append(valueTrailingComments, commentNode)
			p.nextToken()
		}
		value.AddTrailingComments(valueTrailingComments)

		pair := KeyValuePair{
			Key:              key,
			Value:            value,
			LeadingComments:  leadingComments,
			TrailingComments: keyTrailingComments,
		}

		node.Pairs = append(node.Pairs, pair)

		if p.curToken.Type == lexer.TokenComma {
			p.nextToken()
		}
	}

	if p.curToken.Type != lexer.TokenCurlyBraceClose {
		log.Printf("Expected TokenCurlyBraceClose, got %s", p.curToken.Literal)
		return nil
	}

	p.nextToken()
	// 收集对象后的注释
	var trailingComments []*CommentNode
	for p.curToken.Type == lexer.TokenComment {
		commentNode := &CommentNode{
			JSONNode: &JSONNode{},
			Text:     p.curToken.Literal,
		}
		trailingComments = append(trailingComments, commentNode)
		p.nextToken()
	}
	node.AddTrailingComments(trailingComments)

	log.Println("Successfully parsed object")
	return node
}

func (p *Parser) parseArray() *ArrayNode {
	log.Println("Parsing array")
	node := &ArrayNode{
		JSONNode: &JSONNode{},
		Elements: []Node{},
	}

	if p.curToken.Type != lexer.TokenSquareBracketOpen {
		log.Printf("Expected TokenSquareBracketOpen, got %s", p.curToken.Literal)
		return nil
	}

	p.nextToken()

	for p.curToken.Type != lexer.TokenSquareBracketClose && p.curToken.Type != lexer.TokenEOF {
		// 收集前置注释
		var leadingComments []*CommentNode
		for p.curToken.Type == lexer.TokenComment {
			commentNode := &CommentNode{
				JSONNode: &JSONNode{},
				Text:     p.curToken.Literal,
			}
			leadingComments = append(leadingComments, commentNode)
			p.nextToken()
		}

		// 解析元素
		element := p.parseValue()
		if element == nil {
			log.Println("Failed to parse array element")
			return nil
		}
		element.AddLeadingComments(leadingComments)

		// 收集后置注释
		var trailingComments []*CommentNode
		for p.curToken.Type == lexer.TokenComment {
			commentNode := &CommentNode{
				JSONNode: &JSONNode{},
				Text:     p.curToken.Literal,
			}
			trailingComments = append(trailingComments, commentNode)
			p.nextToken()
		}
		element.AddTrailingComments(trailingComments)

		node.Elements = append(node.Elements, element)

		if p.curToken.Type == lexer.TokenComma {
			p.nextToken()
		}
	}

	if p.curToken.Type != lexer.TokenSquareBracketClose {
		log.Printf("Expected TokenSquareBracketClose, got %s", p.curToken.Literal)
		return nil
	}

	p.nextToken()
	// 收集数组后的注释
	var trailingComments []*CommentNode
	for p.curToken.Type == lexer.TokenComment {
		commentNode := &CommentNode{
			JSONNode: &JSONNode{},
			Text:     p.curToken.Literal,
		}
		trailingComments = append(trailingComments, commentNode)
		p.nextToken()
	}
	node.AddTrailingComments(trailingComments)

	log.Println("Successfully parsed array")
	return node
}

func (p *Parser) parseWhitespace() *WhitespaceNode {
	node := &WhitespaceNode{
		JSONNode: &JSONNode{
			Literal: p.curToken.Literal,
		},
		Value: p.curToken.Literal,
	}
	p.nextToken()
	return node
}

func (p *Parser) parseValue() Node {
	switch p.curToken.Type {
	case lexer.TokenString:
		node := &StringNode{
			JSONNode: &JSONNode{
				Literal: p.curToken.Literal,
			},
			Value: p.curToken.Literal,
		}
		p.nextToken()
		return node
	case lexer.TokenNumber:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			log.Printf("Failed to parse number: %s", p.curToken.Literal)
			return nil
		}
		node := &NumberNode{
			JSONNode: &JSONNode{
				Literal: p.curToken.Literal,
			},
			Value: value,
		}
		p.nextToken()
		return node
	case lexer.TokenBoolean:
		value := p.curToken.Literal == "true"
		node := &BoolNode{
			JSONNode: &JSONNode{
				Literal: p.curToken.Literal,
			},
			Value: value,
		}
		p.nextToken()
		return node
	case lexer.TokenNull:
		node := &NullNode{
			JSONNode: &JSONNode{
				Literal: p.curToken.Literal,
			},
		}
		p.nextToken()
		return node
	case lexer.TokenInfinity:
		node := &InfinityNode{
			JSONNode: &JSONNode{
				Literal: p.curToken.Literal,
			},
			Positive: p.curToken.Literal == "Infinity",
		}
		p.nextToken()
		return node
	case lexer.TokenNaN:
		node := &NaNNode{
			JSONNode: &JSONNode{
				Literal: p.curToken.Literal,
			},
		}
		p.nextToken()
		return node
	case lexer.TokenCurlyBraceOpen:
		return p.parseObject()
	case lexer.TokenSquareBracketOpen:
		return p.parseArray()
	default:
		log.Printf("Unexpected token in parseValue: %s", p.curToken.Literal)
		return nil
	}
}
