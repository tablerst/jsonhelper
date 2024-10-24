package parser

import "fmt"

type Node interface {
	TokenLiteral() string
	String() string
	AddLeadingComments(comments []*CommentNode)
	AddTrailingComments(comments []*CommentNode)
	GetLeadingComments() []*CommentNode
	GetTrailingComments() []*CommentNode
}

type JSONNode struct {
	Literal          string
	LeadingComments  []*CommentNode
	TrailingComments []*CommentNode
}

func (jn *JSONNode) TokenLiteral() string {
	return jn.Literal
}

// Implement the methods for JSONNode
func (jn *JSONNode) AddLeadingComments(comments []*CommentNode) {
	jn.LeadingComments = append(jn.LeadingComments, comments...)
}

func (jn *JSONNode) AddTrailingComments(comments []*CommentNode) {
	jn.TrailingComments = append(jn.TrailingComments, comments...)
}

func (jn *JSONNode) GetLeadingComments() []*CommentNode {
	return jn.LeadingComments
}

func (jn *JSONNode) GetTrailingComments() []*CommentNode {
	return jn.TrailingComments
}

type KeyValuePair struct {
	Key              string
	Value            Node
	LeadingComments  []*CommentNode
	TrailingComments []*CommentNode
}

func (kvp *KeyValuePair) String() string {
	result := ""
	// 添加键前的注释
	for _, comment := range kvp.LeadingComments {
		result += comment.String() + "\n"
	}
	result += "\"" + kvp.Key + "\":"
	// 添加键后的注释
	for _, comment := range kvp.TrailingComments {
		result += " " + comment.String()
	}
	result += " " + kvp.Value.String()
	return result
}

type ObjectNode struct {
	*JSONNode
	Pairs []KeyValuePair
}

// Update the String method for ObjectNode
func (on *ObjectNode) String() string {
	result := "{\n"
	for _, pair := range on.Pairs {
		result += pair.String() + ",\n"
	}
	if len(on.Pairs) > 0 {
		result = result[:len(result)-2] // Remove the last comma
	}
	result += "\n}"
	// Add trailing comments after the closing bracket
	for _, comment := range on.GetTrailingComments() {
		result += "\n" + comment.String()
	}
	return result
}

type ArrayNode struct {
	*JSONNode
	Elements []Node
}

func (an *ArrayNode) String() string {
	result := "[\n"
	for _, element := range an.Elements {
		result += element.String() + ",\n"
	}
	if len(an.Elements) > 0 {
		result = result[:len(result)-2]
	}
	result += "\n]"
	// Add trailing comments after the closing bracket
	for _, comment := range an.GetTrailingComments() {
		result += "\n" + comment.String()
	}
	return result
}

type StringNode struct {
	*JSONNode
	Value string
}

func (sn *StringNode) String() string {
	result := ""
	for _, comment := range sn.GetLeadingComments() {
		result += comment.String() + "\n"
	}
	result += "\"" + sn.Value + "\""
	for _, comment := range sn.GetTrailingComments() {
		result += " " + comment.String()
	}
	return result
}

type NumberNode struct {
	*JSONNode
	Value float64
}

func (nn *NumberNode) String() string {
	result := ""
	for _, comment := range nn.GetLeadingComments() {
		result += comment.String() + "\n"
	}
	result += fmt.Sprintf("%v", nn.Value)

	for _, comment := range nn.GetTrailingComments() {
		result += " " + comment.String()
	}
	return result
}

type BoolNode struct {
	*JSONNode
	Value bool
}

func (bn *BoolNode) String() string {
	result := ""
	for _, comment := range bn.GetLeadingComments() {
		result += comment.String() + "\n"
	}
	if bn.Value {
		result += "true"
	} else {
		result += "false"
	}
	for _, comment := range bn.GetTrailingComments() {
		result += " " + comment.String()
	}
	return result
}

type NullNode struct {
	*JSONNode
}

func (nn *NullNode) String() string {
	result := ""
	for _, comment := range nn.GetLeadingComments() {
		result += comment.String() + "\n"
	}
	result += "null"

	for _, comment := range nn.GetTrailingComments() {
		result += " " + comment.String()
	}
	return result
}

type CommentNode struct {
	*JSONNode
	Text             string
	LeadingNewLines  int
	TrailingNewLines int
}

func (cn *CommentNode) String() string {
	return "//" + cn.Text
}

type InfinityNode struct {
	*JSONNode
	Positive bool
}

func (in *InfinityNode) String() string {
	result := ""
	for _, comment := range in.GetLeadingComments() {
		result += comment.String() + "\n"
	}
	if in.Positive {
		result += "Infinity"
	} else {
		result += "-Infinity"
	}

	for _, comment := range in.GetTrailingComments() {
		result += " " + comment.String()
	}
	return result
}

type NaNNode struct {
	*JSONNode
}

func (nn *NaNNode) String() string {
	result := ""
	for _, comment := range nn.GetLeadingComments() {
		result += comment.String() + "\n"
	}
	result += "NaN"

	for _, comment := range nn.GetTrailingComments() {
		result += " " + comment.String()
	}
	return result
}

type WhitespaceNode struct {
	*JSONNode
	Value string
}

func (wn *WhitespaceNode) String() string {
	return wn.Value
}
