package parser

import "fmt"

type Node interface {
	TokenLiteral() string
	String() string
}

type JSONNode struct {
	Literal string
}

func (jn *JSONNode) TokenLiteral() string {
	return jn.Literal
}

type ObjectNode struct {
	JSONNode
	Pairs map[string]Node
}

func (on *ObjectNode) String() string {
	result := "{"
	for key, value := range on.Pairs {
		result += "\"" + key + "\":" + value.String() + ", "
	}
	if len(on.Pairs) > 0 {
		result = result[:len(result)-2]
	}
	result += "}"
	return result
}

type ArrayNode struct {
	JSONNode
	Elements []Node
}

func (an *ArrayNode) String() string {
	result := "["
	for _, element := range an.Elements {
		result += element.String() + ", "
	}
	if len(an.Elements) > 0 {
		result = result[:len(result)-2]
	}
	result += "]"
	return result
}

type StringNode struct {
	JSONNode
	Value string
}

func (sn *StringNode) String() string {
	return "\"" + sn.Value + "\""
}

type NumberNode struct {
	JSONNode
	Value float64
}

func (nn *NumberNode) String() string {
	return fmt.Sprintf("%v", nn.Value)
}

type BoolNode struct {
	JSONNode
	Value bool
}

func (bn *BoolNode) String() string {
	if bn.Value {
		return "true"
	}
	return "false"
}

type NullNode struct {
	JSONNode
}

func (nn *NullNode) String() string {
	return "null"
}

type CommentNode struct {
	JSONNode
	Text string
}

func (cn *CommentNode) String() string {
	return "//" + cn.Text
}

type InfinityNode struct {
	JSONNode
	Positive bool
}

func (in *InfinityNode) String() string {
	if in.Positive {
		return "Infinity"
	}
	return "-Infinity"
}

type NaNNode struct {
	JSONNode
}

func (nn *NaNNode) String() string {
	return "NaN"
}
