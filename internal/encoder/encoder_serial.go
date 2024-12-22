package encoder

import (
	"fmt"
	"github.com/tablerst/jsonhelper/internal/parser"
	"math"
	"strconv"
	"strings"
)

func Encode(node *parser.ASTNode, pretty bool) (string, error) {
	return encodeNode(node, pretty, 0)
}
func encodeNode(node *parser.ASTNode, pretty bool, indent int) (string, error) {
	switch node.Type {
	case parser.ObjectNode:
		return encodeObject(node, pretty, indent)
	case parser.ArrayNode:
		return encodeArray(node, pretty, indent)
	case parser.StringNode:
		return fmt.Sprintf(`"%s"`, escapeString(fmt.Sprintf("%v", node.Value))), nil
	case parser.NumberNode:
		return formatNumber(node.Value)
	case parser.BooleanNode:
		if node.Value.(bool) {
			return "true", nil
		} else {
			return "false", nil
		}
	case parser.NullNode:
		return "null", nil
	case parser.CommentNode:
		// simply ignore comment node at this stage
		return "", nil
	case "property":
		if len(node.Children) > 0 {
			return encodeNode(node.Children[0], pretty, indent)
		}
		return "", nil
	default:
		// handle identifier node in a easy way
		return fmt.Sprintf("%v", node.Value), nil
	}
}

func encodeObject(node *parser.ASTNode, pretty bool, indent int) (string, error) {
	var sb strings.Builder
	sb.WriteString("{")

	children := node.Children
	// property node list
	if pretty && len(children) > 0 {
		sb.WriteString("\n")
	}

	for i, child := range children {
		if child.Type == parser.CommentNode {
			continue
		}
		if child.Type != "property" {
			continue
		}
		// key
		keyStr := escapeString(child.Key)
		if pretty {
			sb.WriteString(strings.Repeat("  ", indent+1))
		}
		sb.WriteString(fmt.Sprintf(`"%s":`, keyStr))
		if pretty {
			sb.WriteString(" ")
		}
		valStr, err := encodeNode(child, pretty, indent+1)
		if err != nil {
			return "", err
		}
		sb.WriteString(valStr)

		if i < len(children)-1 {
			sb.WriteString(",")
		}
		if pretty {
			sb.WriteString("\n")
		}
	}

	if pretty && len(children) > 0 {
		sb.WriteString(strings.Repeat("  ", indent))
	}
	sb.WriteString("}")
	return sb.String(), nil
}

func encodeArray(node *parser.ASTNode, pretty bool, indent int) (string, error) {
	var sb strings.Builder
	sb.WriteString("[")

	if pretty && len(node.Children) > 0 {
		sb.WriteString("\n")
	}

	for i, child := range node.Children {
		if child.Type == parser.CommentNode {
			continue
		}
		if pretty {
			sb.WriteString(strings.Repeat("  ", indent+1))
		}
		valStr, err := encodeNode(child, pretty, indent+1)
		if err != nil {
			return "", err
		}
		sb.WriteString(valStr)
		if i < len(node.Children)-1 {
			sb.WriteString(",")
		}
		if pretty {
			sb.WriteString("\n")
		}
	}

	if pretty && len(node.Children) > 0 {
		sb.WriteString(strings.Repeat("  ", indent))
	}
	sb.WriteString("]")
	return sb.String(), nil
}

func escapeString(s string) string {
	var sb strings.Builder
	for _, ch := range s {
		switch ch {
		case '\\':
			sb.WriteString("\\\\")
		case '"':
			sb.WriteString("\\\"")
		case '\b':
			sb.WriteString("\\b")
		case '\f':
			sb.WriteString("\\f")
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		default:
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}

func formatNumber(val interface{}) (string, error) {
	switch v := val.(type) {
	case float64:
		if math.IsNaN(v) {
			return "NaN", nil
		}
		if math.IsInf(v, 1) {
			return "Infinity", nil
		}
		if math.IsInf(v, -1) {
			return "-Infinity", nil
		}
		return strconv.FormatFloat(v, 'g', -1, 64), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case int:
		return strconv.Itoa(v), nil
	default:
		return "", fmt.Errorf("unsupported number type: %T", val)
	}
}
