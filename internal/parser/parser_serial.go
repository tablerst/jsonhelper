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
		value, err := p.parseObject()
		if err != nil {
			return nil, err
		}
		p.nextToken()
		return value, nil
	case lexer.LBRACKET:
		value, err := p.parseArray()
		if err != nil {
			return nil, err
		}
		p.nextToken()
		return value, nil
	case lexer.STRING:
		value := p.curToken.Value
		p.nextToken()
		return value, nil
	case lexer.NUMBER:
		numberStr := p.curToken.Value
		p.nextToken()
		lowerStr := strings.ToLower(numberStr)
		if lowerStr == "nan" {
			return math.NaN(), nil
		} else if lowerStr == "infinity" || lowerStr == "+infinity" {
			return math.Inf(1), nil
		} else if lowerStr == "-infinity" {
			return math.Inf(-1), nil
		} else if strings.HasPrefix(lowerStr, "0x") || strings.HasPrefix(lowerStr, "+0x") || strings.HasPrefix(lowerStr, "-0x") {
			// 处理十六进制数字
			prefixLen := 2
			if strings.HasPrefix(lowerStr, "+") || strings.HasPrefix(lowerStr, "-") {
				prefixLen = 3
			}
			i, err := strconv.ParseInt(numberStr[prefixLen:], 16, 64)
			if err != nil {
				return nil, err
			}
			if strings.HasPrefix(numberStr, "-") {
				i = -i
			}
			return i, nil
		} else if strings.HasPrefix(numberStr, ".") || strings.HasSuffix(numberStr, ".") || strings.ContainsAny(numberStr, ".eE") {
			// 处理小数
			f, err := strconv.ParseFloat(numberStr, 64)
			if err != nil {
				return nil, err
			}
			return f, nil
		} else {
			i, err := strconv.ParseInt(numberStr, 10, 64)
			if err != nil {
				// 尝试解析为浮点数
				f, err := strconv.ParseFloat(numberStr, 64)
				if err != nil {
					return nil, err
				}
				return f, nil
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
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.curToken)
	}
}

func (p *Parser) parseObject() (JSONObject, error) {
	obj := make(JSONObject)
	p.nextToken()
	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		if p.curToken.Type != lexer.STRING && p.curToken.Type != lexer.IDENTIFIER {
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
		} else if p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
			return nil, fmt.Errorf("expected comma or '}', got %v", p.curToken)
		}
	}
	return obj, nil
}

func (p *Parser) parseArray() (JSONArray, error) {
	array := make(JSONArray, 0)
	p.nextToken()
	for p.curToken.Type != lexer.RBRACKET && p.curToken.Type != lexer.EOF {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		array = append(array, value)
		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
		} else if p.curToken.Type == lexer.RBRACKET {
			break
		} else if p.curToken.Type != lexer.RBRACKET && p.curToken.Type != lexer.EOF {
			return nil, fmt.Errorf("expected comma or ']', got %v", p.curToken)
		}
	}
	return array, nil
}
