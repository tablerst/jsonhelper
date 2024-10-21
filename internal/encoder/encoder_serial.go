package encoder

import (
	"fmt"
	"github.com/tablerst/jsonhelper/internal/parser"
	"math"
	"strconv"
	"strings"
)

func Encode(value parser.JSONValue, pretty bool) (string, error) {
	return encodeValue(value, pretty, 0)
}

func encodeValue(value parser.JSONValue, pretty bool, indent int) (string, error) {
	switch v := value.(type) {
	case parser.JSONObject:
		return encodeObject(v, pretty, indent)
	case parser.JSONArray:
		return encodeArray(v, pretty, indent)
	case string:
		return fmt.Sprintf("\"%s\"", escapeString(v)), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		if math.IsNaN(v) {
			return "NaN", nil
		} else if math.IsInf(v, 1) {
			return "Infinity", nil
		} else if math.IsInf(v, -1) {
			return "-Infinity", nil
		} else {
			return strconv.FormatFloat(v, 'g', -1, 64), nil
		}
	case bool:
		if v {
			return "true", nil
		} else {
			return "false", nil
		}
	case nil:
		return "null", nil
	default:
		return "", fmt.Errorf("unsupported type: %T", value)
	}
}

func encodeObject(obj parser.JSONObject, pretty bool, indent int) (string, error) {
	var sb strings.Builder
	sb.WriteString("{")
	if pretty {
		sb.WriteString("\n")
	}
	i := 0
	for k, v := range obj {
		if pretty {
			sb.WriteString(strings.Repeat("  ", indent+1))
		}
		keyStr := fmt.Sprintf("\"%s\"", escapeString(k))
		sb.WriteString(keyStr)
		sb.WriteString(":")
		if pretty {
			sb.WriteString(" ")
		}
		valStr, err := encodeValue(v, pretty, indent+1)
		if err != nil {
			return "", err
		}
		sb.WriteString(valStr)
		i++
		if i < len(obj) {
			sb.WriteString(",")
			if pretty {
				sb.WriteString("\n")
			}
		}
	}
	if pretty {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("  ", indent))
	}
	sb.WriteString("}")
	return sb.String(), nil
}

func encodeArray(array parser.JSONArray, pretty bool, indent int) (string, error) {
	var sb strings.Builder
	sb.WriteString("[")
	if pretty {
		sb.WriteString("\n")
	}
	for i, v := range array {
		if pretty {
			sb.WriteString(strings.Repeat("  ", indent+1))
		}
		valStr, err := encodeValue(v, pretty, indent+1)
		if err != nil {
			return "", err
		}
		sb.WriteString(valStr)
		if i < len(array)-1 {
			sb.WriteString(",")
			if pretty {
				sb.WriteString("\n")
			}
		}
	}
	if pretty {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("  ", indent))
	}
	sb.WriteString("]")
	return sb.String(), nil
}

func escapeString(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}