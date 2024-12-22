package parser

import (
	"fmt"
	"github.com/tablerst/jsonhelper/internal/lexer"
	"math"
	"strconv"
	"strings"
)

type NodeType string

const (
	ObjectNode     NodeType = "object"
	ArrayNode      NodeType = "array"
	StringNode     NodeType = "string"
	NumberNode     NodeType = "number"
	BooleanNode    NodeType = "boolean"
	NullNode       NodeType = "null"
	CommentNode    NodeType = "comment"
	IdentifierNode NodeType = "identifier"
)

type ASTNode struct {
	Type        NodeType
	Value       interface{}
	Key         string
	Children    []*ASTNode
	Parent      *ASTNode
	StartOffset int
	EndOffset   int
	Line        int
	Column      int
}

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() (*ASTNode, error) {
	return p.parseValue()
}

func (p *Parser) parseValue() (*ASTNode, error) {
	token := p.curToken
	node := &ASTNode{
		StartOffset: token.StartOffset,
		EndOffset:   token.EndOffset,
		Line:        token.Line,
		Column:      token.Column,
	}

	switch token.Type {
	case lexer.LBRACE:
		objectNode, err := p.parseObject()
		if err != nil {
			return nil, err
		}
		p.nextToken()
		return objectNode, nil
	case lexer.LBRACKET:
		arrayNode, err := p.parseArray()
		if err != nil {
			return nil, err
		}
		p.nextToken()
		return arrayNode, nil
	case lexer.STRING:
		node.Type = StringNode
		node.Value = token.Value
		p.nextToken()
		return node, nil
	case lexer.NUMBER:
		node.Type = NumberNode
		val, err := p.parseNumber(token.Value)
		if err != nil {
			return nil, err
		}
		node.Value = val
		p.nextToken()
		return node, nil
	case lexer.BOOLEAN:
		node.Type = BooleanNode
		node.Value = strings.ToLower(p.curToken.Value) == "true"
		p.nextToken()
		return node, nil
	case lexer.NULL:
		node.Type = NullNode
		node.Value = nil
		p.nextToken()
		return node, nil
	case lexer.COMMENT:
		// TODO CommentNode waited to be implemented
		node.Type = CommentNode
		node.Value = token.Value
		p.nextToken()
		return node, nil
	case lexer.IDENTIFIER:
		node.Type = IdentifierNode
		node.Value = token.Value
		p.nextToken()
		return node, nil
	case lexer.EOF:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", token)
	}
}

func (p *Parser) parseObject() (*ASTNode, error) {
	start := p.curToken.StartOffset
	line := p.curToken.Line
	column := p.curToken.Column

	objNode := &ASTNode{
		Type:        ObjectNode,
		StartOffset: start,
		Line:        line,
		Column:      column,
		Children:    make([]*ASTNode, 0),
	}

	p.nextToken()
	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		var keyNode *ASTNode
		if p.curToken.Type == lexer.STRING || p.curToken.Type == lexer.IDENTIFIER {
			keyNode = &ASTNode{
				Type:        StringNode,
				Value:       p.curToken.Value,
				StartOffset: p.curToken.StartOffset,
				EndOffset:   p.curToken.EndOffset,
				Line:        p.curToken.Line,
				Column:      p.curToken.Column,
			}
			p.nextToken()
		} else if p.curToken.Type == lexer.COMMENT {
			// TODO CommentNode waited to be implemented
			p.nextToken()
			continue
		} else {
			return nil, fmt.Errorf("expected property  key, got %v", p.curToken)
		}

		if p.curToken.Type != lexer.COLON {
			return nil, fmt.Errorf("expected colon after key, got %v", p.curToken)
		}
		p.nextToken()

		valNode, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		propertyNode := &ASTNode{
			Type:        "property",
			Key:         keyNode.Value.(string),
			StartOffset: keyNode.StartOffset,
			// 结束位置先暂时定在 value 结束
			EndOffset: valNode.EndOffset,
			Children:  []*ASTNode{valNode},
		}
		valNode.Parent = propertyNode
		objNode.Children = append(objNode.Children, propertyNode)

		// May be a trailing comma
		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
			if p.curToken.Type == lexer.RBRACE {
				// trailing comma
				break
			}
		} else if p.curToken.Type == lexer.RBRACE {
			break
		} else if p.curToken.Type == lexer.EOF {
			break
		} else if p.curToken.Type == lexer.COMMENT {
			p.nextToken()
		} else {
			return nil, fmt.Errorf("expected comma or '}', got %v", p.curToken)
		}
	}
	if p.curToken.Type == lexer.RBRACE {
		objNode.EndOffset = p.curToken.EndOffset
		p.nextToken()
	}
	return objNode, nil
}

func (p *Parser) parseArray() (*ASTNode, error) {
	start := p.curToken.StartOffset
	line := p.curToken.Line
	col := p.curToken.Column

	arrayNode := &ASTNode{
		Type:        ArrayNode,
		StartOffset: start,
		Line:        line,
		Column:      col,
		Children:    []*ASTNode{},
	}
	p.nextToken()

	for p.curToken.Type != lexer.RBRACKET && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.COMMENT {
			p.nextToken()
			continue
		}
		valNode, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		if valNode != nil {
			valNode.Parent = arrayNode
			arrayNode.Children = append(arrayNode.Children, valNode)
		}

		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
			if p.curToken.Type == lexer.RBRACKET {
				break
			}
		} else if p.curToken.Type == lexer.RBRACKET {
			break
		} else if p.curToken.Type == lexer.EOF {
			break
		} else {
			return nil, fmt.Errorf("expected comma or ']', got %v", p.curToken)
		}
	}

	if p.curToken.Type == lexer.RBRACKET {
		arrayNode.EndOffset = p.curToken.EndOffset
		p.nextToken()
	}

	return arrayNode, nil
}

func (p *Parser) parseNumber(s string) (interface{}, error) {
	lower := strings.ToLower(s)
	switch lower {
	case "nan":
		return math.NaN(), nil
	case "infinity", "+infinity":
		return math.Inf(1), nil
	case "-infinity":
		return math.Inf(-1), nil
	}

	if strings.HasPrefix(lower, "0x") || strings.HasPrefix(lower, "+0x") || strings.HasPrefix(lower, "-0x") {
		prefixLen := 2
		if strings.HasPrefix(lower, "+") || strings.HasPrefix(lower, "-") {
			prefixLen = 3
		}
		i, err := strconv.ParseInt(s[prefixLen:], 16, 64)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(s, "-") {
			i = -i
		}
		return i, nil
	}

	if strings.ContainsAny(s, ".eE") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		f, err2 := strconv.ParseFloat(s, 64)
		if err2 != nil {
			return nil, err2
		}
		return f, nil
	}
	return i, nil
}
