package test

import (
	"github.com/tablerst/jsonhelper/internal/lexer"
	"github.com/tablerst/jsonhelper/internal/parser"
	"testing"
)

func TestParseJSON5(t *testing.T) {
	input := `{
		// This is a comment
		name: 'ChatGPT',
		age: 3,
		is_active: true,
		address: {
			city: "San Francisco",
			zip: 94105,
		},
		values: [1, 2, 3, "end"]
	}`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)

	result := p.Parser()

	if len(result.Errors) > 0 {
		t.Fatalf("Parser encountered errors: %v", result.Errors)
	}

	if result.Root == nil {
		t.Fatalf("Expected AST root node, got nil")
	}

	objectNode, ok := result.Root.(*parser.ObjectNode)
	if !ok {
		t.Fatalf("Expected ObjectNode, got %T", result.Root)
	}

	// Check name field
	if nameNode, exists := objectNode.Pairs["name"]; !exists {
		t.Errorf("Expected key 'name' to exist")
	} else if nameNode == nil {
		t.Errorf("Node for key 'name' is nil")
	} else {
		nameStrNode, ok := nameNode.(*parser.StringNode)
		if !ok {
			t.Errorf("Expected StringNode for key 'name', got %T", nameNode)
		}
		if nameStrNode.Value != "ChatGPT" {
			t.Errorf("Expected value 'ChatGPT' for key 'name', got %s", nameStrNode.Value)
		}
	}

	// Check age field
	if ageNode, exists := objectNode.Pairs["age"]; !exists {
		t.Errorf("Expected key 'age' to exist")
	} else {
		ageNumNode, ok := ageNode.(*parser.NumberNode)
		if !ok {
			t.Errorf("Expected NumberNode for key 'age', got %T", ageNode)
		}
		if ageNumNode.Value != 3 {
			t.Errorf("Expected value 3 for key 'age', got %f", ageNumNode.Value)
		}
	}

	// Check is_active field
	if isActiveNode, exists := objectNode.Pairs["is_active"]; !exists {
		t.Errorf("Expected key 'is_active' to exist")
	} else {
		isActiveBoolNode, ok := isActiveNode.(*parser.BoolNode)
		if !ok {
			t.Errorf("Expected BoolNode for key 'is_active', got %T", isActiveNode)
		}
		if isActiveBoolNode.Value != true {
			t.Errorf("Expected value true for key 'is_active', got %v", isActiveBoolNode.Value)
		}
	}

	// Check address field
	if addressNode, exists := objectNode.Pairs["address"]; !exists {
		t.Errorf("Expected key 'address' to exist")
	} else {
		addressObjectNode, ok := addressNode.(*parser.ObjectNode)
		if !ok {
			t.Fatalf("Expected ObjectNode for key 'address', got %T", addressNode)
		}

		// Check city field inside address
		if cityNode, exists := addressObjectNode.Pairs["city"]; !exists {
			t.Errorf("Expected key 'city' to exist in 'address'")
		} else {
			cityStrNode, ok := cityNode.(*parser.StringNode)
			if !ok || cityStrNode.Value != "San Francisco" {
				t.Errorf("Expected StringNode with value 'San Francisco' for key 'city', got %v", cityNode)
			}
		}

		// Check zip field inside address
		if zipNode, exists := addressObjectNode.Pairs["zip"]; !exists {
			t.Errorf("Expected key 'zip' to exist in 'address'")
		} else {
			zipNumNode, ok := zipNode.(*parser.NumberNode)
			if !ok || zipNumNode.Value != 94105 {
				t.Errorf("Expected NumberNode with value 94105 for key 'zip', got %v", zipNode)
			}
		}
	}

	// Check values array field
	if valuesNode, exists := objectNode.Pairs["values"]; !exists {
		t.Errorf("Expected key 'values' to exist")
	} else {
		valuesArrayNode, ok := valuesNode.(*parser.ArrayNode)
		if !ok {
			t.Fatalf("Expected ArrayNode for key 'values', got %T", valuesNode)
		}

		if len(valuesArrayNode.Elements) != 4 {
			t.Errorf("Expected 4 elements in 'values' array, got %d", len(valuesArrayNode.Elements))
		}

		for i, expected := range []interface{}{1.0, 2.0, 3.0, "end"} {
			element := valuesArrayNode.Elements[i]
			switch v := expected.(type) {
			case float64:
				numNode, ok := element.(*parser.NumberNode)
				if !ok || numNode.Value != v {
					t.Errorf("Expected NumberNode with value %v, got %v", v, element)
				}
			case string:
				strNode, ok := element.(*parser.StringNode)
				if !ok || strNode.Value != v {
					t.Errorf("Expected StringNode with value '%v', got %v", v, element)
				}
			}
		}
	}
}
