package test

import (
	"github.com/tablerst/jsonhelper/internal/lexer"
	"github.com/tablerst/jsonhelper/internal/parser"
	"strings"
	"testing"
)

// Helper function to trim braces and whitespace for comparison.
func trimBracesAndWhitespace(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "{") {
		s = strings.TrimPrefix(s, "{")
	}
	if strings.HasSuffix(s, "}") {
		s = strings.TrimSuffix(s, "}")
	}
	s = strings.TrimSpace(s)
	return s
}

func TestParserWithComments(t *testing.T) {
	input := `
    {
        // This is a comment for key1
        "key1": "value1", // End of key1
        "key2": 123,
        // Comment after key2
        "key3": [
            // Comment for first element
            "elem1",
            "elem2" // End of elem2
        ] // End of key3 array
    }
    `

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	result := p.Parser()

	if len(result.Errors) != 0 {
		t.Fatalf("Parser has errors: %v", result.Errors)
	}

	expected := `// This is a comment for key1
"key1": "value1",// End of key1
"key2": 123,
// Comment after key2
"key3": [
    // Comment for first element
    "elem1",
    "elem2" // End of elem2
] // End of key3 array`

	output := result.Root.String()

	// Remove outer braces and trim whitespace
	output = trimBracesAndWhitespace(output)
	expected = trimBracesAndWhitespace(expected)
	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

// JSON5 Test Case
func TestParserWithJSON5Comments(t *testing.T) {
	input := `
    // Root level comment
    {
        /* Comment for key1 */
        "key1": "value1", // End of key1
        "key2": 123, // Number value
        "key3": {
            // Nested object comment
            "nestedKey1": true,
            "nestedKey2": null // Null value
        },
        "key4": [
            // Array comment
            "elem1",
            "elem2", // Second element
            // Last element comment
            "elem3"
        ],
        "key5": Infinity, // Infinity value
        "key6": -Infinity, // Negative Infinity
        "key7": NaN // Not a Number
    }
    `

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	result := p.Parser()

	if len(result.Errors) != 0 {
		t.Fatalf("Parser has errors: %v", result.Errors)
	}

	expected := `// Root level comment
/* Comment for key1 */
"key1": "value1", // End of key1
"key2": 123, // Number value
"key3": {
    // Nested object comment
    "nestedKey1": true,
    "nestedKey2": null // Null value
},
"key4": [
    // Array comment
    "elem1",
    "elem2", // Second element
    // Last element comment
    "elem3"
],
"key5": Infinity, // Infinity value
"key6": -Infinity, // Negative Infinity
"key7": NaN // Not a Number`

	output := result.Root.String()

	// Remove outer braces and trim whitespace
	output = trimBracesAndWhitespace(output)
	expected = trimBracesAndWhitespace(expected)
	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}
