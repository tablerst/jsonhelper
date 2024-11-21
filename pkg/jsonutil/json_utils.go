package jsonutil

import (
	"github.com/tablerst/jsonhelper/internal/encoder"
	"github.com/tablerst/jsonhelper/internal/lexer"
	"github.com/tablerst/jsonhelper/internal/parser"
)

func Parse(jsonStr string) (interface{}, error) {
	l := lexer.New(jsonStr)
	p := parser.New(l)

	result, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Encode(parseData interface{}, pretty bool) (string, error) {
	output, err := encoder.Encode(parseData, pretty)
	if err != nil {
		return "", err
	}

	return output, nil
}
