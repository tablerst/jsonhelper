package parser

import (
	"fmt"
	"github.com/tablerst/jsonhelper/internal/lexer"
	"math"
	"strconv"
	"strings"
)

type JSONValue interface{}

type JSONObject map[string]JSONValue
type JSONArray []JSONValue

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

func (p *Parser) Parse() (JSONValue, error) {
	return p.parseValue()
}

func (p *Parser) parseValue() (JSONValue, error) {
	switch p.curToken.Type {
	case lexer.LBRACE:
		return p.parseObject()
	case lexer.LBRACKET:
		return p.parseArray()
	case lexer.STRING:
		value := p.curToken.Value
		p.nextToken()
		return value, nil
	case lexer.NUMBER:
		numberStr := p.curToken.Value
		p.nextToken()
		// 处理特殊值 NaN 和 Infinity
		lowerStr := strings.ToLower(numberStr)
		if lowerStr == "nan" {
			return math.NaN(), nil
		} else if lowerStr == "infinity" || lowerStr == "+infinity" {
			return math.Inf(1), nil
		} else if lowerStr == "-infinity" {
			return math.Inf(-1), nil
		}

		if strings.ContainsAny(numberStr, ".eE") {
			f, err := strconv.ParseFloat(numberStr, 64)
			if err != nil {
				return nil, err
			}
			return f, nil
		} else {
			i, err := strconv.ParseInt(numberStr, 10, 64)
			if err != nil {
				return nil, err
			}
			return i, nil
		}
	case lexer.BOOLEAN:
		value := strings.ToLower(p.curToken.Value) == "true"
		p.nextToken()
		return value, nil
	case lexer.NULL:
		p.nextToken()
		return nil, nil
	case lexer.NAN:
		p.nextToken()
		return "NaN", nil
	case lexer.INFINITY:
		p.nextToken()
		return "Infinity", nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.curToken)
	}
}

func (p *Parser) parseObject() (JSONObject, error) {
	obj := make(JSONObject)
	p.nextToken()
	for p.curToken.Type != lexer.RBRACE {
		if p.curToken.Type != lexer.STRING {
			return nil, fmt.Errorf("expected string key, got %v", p.curToken)
		}
		key := p.curToken.Value
		p.nextToken()
		if p.curToken.Type != lexer.COLON {
			return nil, fmt.Errorf("expected colon after key, got %v", p.curToken)
		}
		p.nextToken()
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = value
		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
		} else if p.curToken.Type == lexer.RBRACE {
			break
		} else {
			return nil, fmt.Errorf("expected comma or '}', got %v", p.curToken)
		}
	}
	p.nextToken()
	return obj, nil
}

func (p *Parser) parseArray() (JSONArray, error) {
	array := make(JSONArray, 0)
	p.nextToken()
	for p.curToken.Type != lexer.RBRACKET {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		array = append(array, value)
		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
		} else if p.curToken.Type == lexer.RBRACKET {
			break
		} else {
			return nil, fmt.Errorf("expected comma or ']', got %v", p.curToken)
		}
	}
	p.nextToken()
	return array, nil
}
