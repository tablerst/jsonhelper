package test

import (
	"fmt"
	"github.com/tablerst/jsonhelper/internal/encoder"
	"github.com/tablerst/jsonhelper/internal/lexer"
	"github.com/tablerst/jsonhelper/internal/parser"
	"testing"
)

func TestEncoder_Encode(t *testing.T) {
	jsonStr := `
{
  // comments
  unquoted: 'and you can quote me on that',
  singleQuotes: 'I can use "double quotes" here',
  lineBreaks: "Look, Mom! \
No \\n's!",
  hexadecimal: 0xdecaf,
  leadingDecimalPoint: .8675309, andTrailing: 8675309.,
  positiveSign: +1,
  trailingComma: 'in objects', andIn: ['arrays',],
  "backwardsCompatible": "with JSON",
}
`

	// 创建词法分析器
	l := lexer.New(jsonStr)

	// 创建语法分析器
	p := parser.New(l)

	// 解析 JSON
	result, err := p.Parse()
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// 编码 JSON，使用美化输出
	output, err := encoder.Encode(result, true)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("Parsed and encoded JSON:")
	fmt.Println(output)
}
