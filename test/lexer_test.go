package test

import (
	"github.com/tablerst/jsonhelper/internal/lexer"
	"testing"
)

func TestLexer_JSON(t *testing.T) {
	input := `{
		"name": "John Doe",
		"age": 30,
		"isEmployed": true,
		"address": {
			"street": "123 Main St",
			"city": "Anytown"
		},
		"hobbies": ["reading", "gaming", "hiking"],
		"spouse": null
	}`

	l := lexer.NewLexer(input)

	expectedTokens := []lexer.Token{
		{Type: lexer.TokenCurlyBraceOpen, Literal: "{", Line: 1, Column: 1},
		{Type: lexer.TokenString, Literal: "name", Line: 2, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 2, Column: 9},
		{Type: lexer.TokenString, Literal: "John Doe", Line: 2, Column: 11},
		{Type: lexer.TokenComma, Literal: ",", Line: 2, Column: 21},
		{Type: lexer.TokenString, Literal: "age", Line: 3, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 3, Column: 7},
		{Type: lexer.TokenNumber, Literal: "30", Line: 3, Column: 9},
		{Type: lexer.TokenComma, Literal: ",", Line: 3, Column: 11},
		{Type: lexer.TokenString, Literal: "isEmployed", Line: 4, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 4, Column: 14},
		{Type: lexer.TokenBoolean, Literal: "true", Line: 4, Column: 16},
		{Type: lexer.TokenComma, Literal: ",", Line: 4, Column: 20},
		{Type: lexer.TokenString, Literal: "address", Line: 5, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 5, Column: 11},
		{Type: lexer.TokenCurlyBraceOpen, Literal: "{", Line: 5, Column: 13},
		{Type: lexer.TokenString, Literal: "street", Line: 6, Column: 5},
		{Type: lexer.TokenColon, Literal: ":", Line: 6, Column: 13},
		{Type: lexer.TokenString, Literal: "123 Main St", Line: 6, Column: 15},
		{Type: lexer.TokenComma, Literal: ",", Line: 6, Column: 26},
		{Type: lexer.TokenString, Literal: "city", Line: 7, Column: 5},
		{Type: lexer.TokenColon, Literal: ":", Line: 7, Column: 10},
		{Type: lexer.TokenString, Literal: "Anytown", Line: 7, Column: 12},
		{Type: lexer.TokenCurlyBraceClose, Literal: "}", Line: 8, Column: 3},
		{Type: lexer.TokenComma, Literal: ",", Line: 8, Column: 4},
		{Type: lexer.TokenString, Literal: "hobbies", Line: 9, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 9, Column: 11},
		{Type: lexer.TokenSquareBracketOpen, Literal: "[", Line: 9, Column: 13},
		{Type: lexer.TokenString, Literal: "reading", Line: 9, Column: 14},
		{Type: lexer.TokenComma, Literal: ",", Line: 9, Column: 23},
		{Type: lexer.TokenString, Literal: "gaming", Line: 9, Column: 25},
		{Type: lexer.TokenComma, Literal: ",", Line: 9, Column: 32},
		{Type: lexer.TokenString, Literal: "hiking", Line: 9, Column: 34},
		{Type: lexer.TokenSquareBracketClose, Literal: "]", Line: 9, Column: 42},
		{Type: lexer.TokenComma, Literal: ",", Line: 9, Column: 43},
		{Type: lexer.TokenString, Literal: "spouse", Line: 10, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 10, Column: 10},
		{Type: lexer.TokenNull, Literal: "null", Line: 10, Column: 12},
		{Type: lexer.TokenCurlyBraceClose, Literal: "}", Line: 11, Column: 1},
		{Type: lexer.TokenEOF, Literal: "", Line: 11, Column: 2},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != expected.Type || tok.Literal != expected.Literal {
			t.Errorf("Test JSON - Token %d: expected %v (%s), got %v (%s)", i, expected.Type, expected.Literal, tok.Type, tok.Literal)
		}
	}
}

// TestJSON5Lexer tests the lexing of JSON5 features
func TestLexer_JSON5(t *testing.T) {
	input := `{
		// User information
		name: 'Jane Doe', // Trailing comma allowed
		age: -25,
		isEmployed: false,
		address: {
			street: "456 Elm St",
			city: 'Othertown',
			coordinates: { lat: 40.7128, lng: -74.0060, }, /* Multi-line
				comment */
		},
		hobbies: ['painting', 'cycling',],
		website: undefined, // JSON5 allows undefined
		infinityValue: Infinity,
		notANumber: NaN,
	}`

	l := lexer.NewLexer(input)

	expectedTokens := []lexer.Token{
		{Type: lexer.TokenCurlyBraceOpen, Literal: "{", Line: 1, Column: 1},
		{Type: lexer.TokenComment, Literal: " User information", Line: 2, Column: 3},
		{Type: lexer.TokenString, Literal: "name", Line: 3, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 3, Column: 7},
		{Type: lexer.TokenString, Literal: "Jane Doe", Line: 3, Column: 9},
		{Type: lexer.TokenComma, Literal: ",", Line: 3, Column: 19},
		{Type: lexer.TokenComment, Literal: " Trailing comma allowed", Line: 3, Column: 20},
		{Type: lexer.TokenString, Literal: "age", Line: 4, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 4, Column: 7},
		{Type: lexer.TokenNumber, Literal: "-25", Line: 4, Column: 9},
		{Type: lexer.TokenComma, Literal: ",", Line: 4, Column: 12},
		{Type: lexer.TokenString, Literal: "isEmployed", Line: 5, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 5, Column: 14},
		{Type: lexer.TokenBoolean, Literal: "false", Line: 5, Column: 16},
		{Type: lexer.TokenComma, Literal: ",", Line: 5, Column: 21},
		{Type: lexer.TokenString, Literal: "address", Line: 6, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 6, Column: 11},
		{Type: lexer.TokenCurlyBraceOpen, Literal: "{", Line: 6, Column: 13},
		{Type: lexer.TokenString, Literal: "street", Line: 7, Column: 5},
		{Type: lexer.TokenColon, Literal: ":", Line: 7, Column: 13},
		{Type: lexer.TokenString, Literal: "456 Elm St", Line: 7, Column: 15},
		{Type: lexer.TokenComma, Literal: ",", Line: 7, Column: 25},
		{Type: lexer.TokenString, Literal: "city", Line: 8, Column: 5},
		{Type: lexer.TokenColon, Literal: ":", Line: 8, Column: 10},
		{Type: lexer.TokenString, Literal: "Othertown", Line: 8, Column: 12},
		{Type: lexer.TokenComma, Literal: ",", Line: 8, Column: 22},
		{Type: lexer.TokenString, Literal: "coordinates", Line: 9, Column: 5},
		{Type: lexer.TokenColon, Literal: ":", Line: 9, Column: 16},
		{Type: lexer.TokenCurlyBraceOpen, Literal: "{", Line: 9, Column: 18},
		{Type: lexer.TokenString, Literal: "lat", Line: 9, Column: 20},
		{Type: lexer.TokenColon, Literal: ":", Line: 9, Column: 23},
		{Type: lexer.TokenNumber, Literal: "40.7128", Line: 9, Column: 25},
		{Type: lexer.TokenComma, Literal: ",", Line: 9, Column: 31},
		{Type: lexer.TokenString, Literal: "lng", Line: 9, Column: 33},
		{Type: lexer.TokenColon, Literal: ":", Line: 9, Column: 36},
		{Type: lexer.TokenNumber, Literal: "-74.0060", Line: 9, Column: 38},
		{Type: lexer.TokenComma, Literal: ",", Line: 9, Column: 46},
		{Type: lexer.TokenCurlyBraceClose, Literal: "}", Line: 9, Column: 47},
		{Type: lexer.TokenComma, Literal: ",", Line: 9, Column: 48},
		{Type: lexer.TokenComment, Literal: " Multi-line\n\t\t\t\tcomment ", Line: 9, Column: 49},
		{Type: lexer.TokenCurlyBraceClose, Literal: "}", Line: 10, Column: 3},
		{Type: lexer.TokenComma, Literal: ",", Line: 10, Column: 4},
		{Type: lexer.TokenString, Literal: "hobbies", Line: 11, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 11, Column: 11},
		{Type: lexer.TokenSquareBracketOpen, Literal: "[", Line: 11, Column: 13},
		{Type: lexer.TokenString, Literal: "painting", Line: 11, Column: 14},
		{Type: lexer.TokenComma, Literal: ",", Line: 11, Column: 23},
		{Type: lexer.TokenString, Literal: "cycling", Line: 11, Column: 25},
		{Type: lexer.TokenComma, Literal: ",", Line: 11, Column: 32},
		{Type: lexer.TokenSquareBracketClose, Literal: "]", Line: 11, Column: 33},
		{Type: lexer.TokenComma, Literal: ",", Line: 11, Column: 34},
		{Type: lexer.TokenString, Literal: "website", Line: 12, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 12, Column: 11},
		{Type: lexer.TokenString, Literal: "undefined", Line: 12, Column: 13},
		{Type: lexer.TokenComma, Literal: ",", Line: 12, Column: 22},
		{Type: lexer.TokenComment, Literal: " JSON5 allows undefined", Line: 12, Column: 23},
		{Type: lexer.TokenString, Literal: "infinityValue", Line: 13, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 13, Column: 17},
		{Type: lexer.TokenInfinity, Literal: "Infinity", Line: 13, Column: 19},
		{Type: lexer.TokenComma, Literal: ",", Line: 13, Column: 27},
		{Type: lexer.TokenString, Literal: "notANumber", Line: 14, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 14, Column: 14},
		{Type: lexer.TokenNaN, Literal: "NaN", Line: 14, Column: 16},
		{Type: lexer.TokenComma, Literal: ",", Line: 14, Column: 19},
		{Type: lexer.TokenCurlyBraceClose, Literal: "}", Line: 15, Column: 1},
		{Type: lexer.TokenEOF, Literal: "", Line: 15, Column: 2},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != expected.Type || tok.Literal != expected.Literal {
			t.Errorf("Test JSON5 - Token %d: expected %v (%s), got %v (%s)", i, expected.Type, expected.Literal, tok.Type, tok.Literal)
		}
	}
}

// TestJSONCLexer tests the lexing of JSONC features
func TestLexer_JSONC(t *testing.T) {
	input := `{
	// This is a single-line comment
	"name": "Alice",
	"age": 28, /* This is a 
			multi-line comment */
	"skills": ["Go", "JavaScript", "Python"],
	// Trailing comma is allowed in JSONC
}`

	l := lexer.NewLexer(input)

	expectedTokens := []lexer.Token{
		{Type: lexer.TokenCurlyBraceOpen, Literal: "{", Line: 1, Column: 1},
		{Type: lexer.TokenComment, Literal: " This is a single-line comment", Line: 2, Column: 3},
		{Type: lexer.TokenString, Literal: "name", Line: 3, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 3, Column: 9},
		{Type: lexer.TokenString, Literal: "Alice", Line: 3, Column: 11},
		{Type: lexer.TokenComma, Literal: ",", Line: 3, Column: 18},
		{Type: lexer.TokenString, Literal: "age", Line: 4, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 4, Column: 7},
		{Type: lexer.TokenNumber, Literal: "28", Line: 4, Column: 9},
		{Type: lexer.TokenComma, Literal: ",", Line: 4, Column: 11},
		{Type: lexer.TokenComment, Literal: " This is a \n\t\t\tmulti-line comment ", Line: 4, Column: 12},
		{Type: lexer.TokenString, Literal: "skills", Line: 6, Column: 3},
		{Type: lexer.TokenColon, Literal: ":", Line: 6, Column: 11},
		{Type: lexer.TokenSquareBracketOpen, Literal: "[", Line: 6, Column: 13},
		{Type: lexer.TokenString, Literal: "Go", Line: 6, Column: 14},
		{Type: lexer.TokenComma, Literal: ",", Line: 6, Column: 17},
		{Type: lexer.TokenString, Literal: "JavaScript", Line: 6, Column: 19},
		{Type: lexer.TokenComma, Literal: ",", Line: 6, Column: 30},
		{Type: lexer.TokenString, Literal: "Python", Line: 6, Column: 32},
		{Type: lexer.TokenSquareBracketClose, Literal: "]", Line: 6, Column: 39},
		{Type: lexer.TokenComma, Literal: ",", Line: 6, Column: 40},
		{Type: lexer.TokenComment, Literal: " Trailing comma is allowed in JSONC", Line: 7, Column: 3},
		{Type: lexer.TokenCurlyBraceClose, Literal: "}", Line: 7, Column: 4},
		{Type: lexer.TokenEOF, Literal: "", Line: 7, Column: 5},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != expected.Type || tok.Literal != expected.Literal {
			t.Errorf("Test JSONC - Token %d: expected %v (%s), got %v (%s)", i, expected.Type, expected.Literal, tok.Type, tok.Literal)
		}
	}
}

// TestLexerSkipsComments tests that the lexer skips comments todo?
func TestLexerSkipsComments(t *testing.T) {
	input := `{
        // This is a comment
        name: 'ChatGPT',
        /* Multi-line
           comment */
        age: 3
    }`

	l := lexer.NewLexer(input)
	tokens := []lexer.Token{}
	for tok := l.NextToken(); tok.Type != lexer.TokenEOF; tok = l.NextToken() {
		tokens = append(tokens, tok)
	}

	expectedTokens := []lexer.Token{
		{Type: lexer.TokenCurlyBraceOpen, Literal: "{"},
		{Type: lexer.TokenInfinity, Literal: "name"},
		{Type: lexer.TokenColon, Literal: ":"},
		{Type: lexer.TokenString, Literal: "ChatGPT"},
		{Type: lexer.TokenComma, Literal: ","},
		{Type: lexer.TokenInfinity, Literal: "age"},
		{Type: lexer.TokenColon, Literal: ":"},
		{Type: lexer.TokenNumber, Literal: "3"},
		{Type: lexer.TokenCurlyBraceClose, Literal: "}"},
	}

	if len(tokens) != len(expectedTokens) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTokens), len(tokens))
	}

	for i, tok := range tokens {
		if tok.Type != expectedTokens[i].Type || tok.Literal != expectedTokens[i].Literal {
			t.Errorf("Token %d - expected (%v, %q), got (%v, %q)", i, expectedTokens[i].Type, expectedTokens[i].Literal, tok.Type, tok.Literal)
		}
	}
}
