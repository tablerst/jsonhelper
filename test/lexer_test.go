package test

import (
	"github.com/tablerst/jsonhelper/internal/lexer"
	"testing"
)

func TestEmptyInput(t *testing.T) {
	input := ""
	l := lexer.NewLexer(input)
	tok := l.NextToken()
	if tok.Type != lexer.TokenEOF {
		t.Fatalf("Expected EOF token, got %s", tok.Type)
	}
}

func TestSingleTokens(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{" ", lexer.TokenWhitespace, " "},
		{"\t", lexer.TokenWhitespace, "\t"},
		{"\n", lexer.TokenNewLine, "\n"},
		{"-", lexer.TokenMinus, "-"},
		{"{", lexer.TokenCurlyBraceOpen, "{"},
		{"}", lexer.TokenCurlyBraceClose, "}"},
		{"[", lexer.TokenSquareBracketOpen, "["},
		{"]", lexer.TokenSquareBracketClose, "]"},
		{",", lexer.TokenComma, ","},
		{":", lexer.TokenColon, ":"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("Input %q - Expected token type %s, got %s", tt.input, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("Input %q - Expected literal %s, got %s", tt.input, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestStrings(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`"hello\nworld"`, "hello\nworld"},
		{`"hello\u0041"`, "helloA"},
		{`"hello\u{1F600}"`, "hello😀"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		tok := l.NextToken()
		if tok.Type != lexer.TokenString {
			t.Errorf("Input %q - Expected token type STRING, got %s", tt.input, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("Input %q - Expected literal %s, got %s", tt.input, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNumbers(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{"123", "123"},
		{"0", "0"},
		{"0.123", "0.123"},
		{"123.456", "123.456"},
		{"1e10", "1e10"},
		{"1E10", "1E10"},
		{"1e-10", "1e-10"},
		{"0x1A3F", "0x1A3F"},
		{"0b1010", "0b1010"},
		{"0o755", "0o755"},
		{"-42", "-42"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		tok := l.NextToken()
		if tok.Type != lexer.TokenNumber {
			t.Errorf("Input %q - Expected token type NUMBER, got %s", tt.input, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("Input %q - Expected literal %s, got %s", tt.input, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIdentifiers(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{"true", lexer.TokenBoolean, "true"},
		{"false", lexer.TokenBoolean, "false"},
		{"null", lexer.TokenNull, "null"},
		{"Infinity", lexer.TokenInfinity, "Infinity"},
		{"NaN", lexer.TokenNaN, "NaN"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("Input %q - Expected token type %s, got %s", tt.input, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("Input %q - Expected literal %s, got %s", tt.input, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestComments(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{"// This is a comment\n", "// This is a comment"},
		{"/* This is a \n multi-line comment */", "/* This is a \n multi-line comment */"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		tok := l.NextToken()
		if tok.Type != lexer.TokenComment {
			t.Errorf("Input %q - Expected token type COMMENT, got %s", tt.input, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("Input %q - Expected literal %s, got %s", tt.input, tt.expectedLiteral, tok.Literal)
		}
	}
}

// test/lexer_test.go

func TestInvalidInputs(t *testing.T) {
	tests := []string{
		"\"unclosed string",
		"0xGHIJ",          // Invalid hex
		"0b102",           // Invalid binary
		"0o789",           // Invalid octal
		"\"\\u123\"",      // Invalid unicode escape
		"\"\\u{110000}\"", // Unicode code point out of range
		"/* unclosed comment",
	}

	for _, input := range tests {
		l := lexer.NewLexer(input)
		tok := l.NextToken()
		if tok.Type != lexer.TokenError {
			t.Errorf("Input %q - Expected TokenError due to invalid input, got %s", input, tok.Type)
		}
	}
}

func TestComplexInput(t *testing.T) {
	input := `{
        "string": "value",
        "number": 12345,
        "object": {
            "bool": true,
            "null": null,
            "array": [1, 2, 3]
        },
        "comment": /* multi-line
                      comment */
                    // single-line comment
                    "end": Infinity
    }`

	expectedTokens := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.TokenCurlyBraceOpen, "{"},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "        "},
		{lexer.TokenString, "string"},
		{lexer.TokenColon, ":"},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenString, "value"},
		{lexer.TokenComma, ","},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "        "},
		{lexer.TokenString, "number"},
		{lexer.TokenColon, ":"},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenNumber, "12345"},
		{lexer.TokenComma, ","},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "        "},
		{lexer.TokenString, "object"},
		{lexer.TokenColon, ":"},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenCurlyBraceOpen, "{"},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "            "},
		{lexer.TokenString, "bool"},
		{lexer.TokenColon, ":"},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenBoolean, "true"},
		{lexer.TokenComma, ","},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "            "},
		{lexer.TokenString, "null"},
		{lexer.TokenColon, ":"},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenNull, "null"},
		{lexer.TokenComma, ","},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "            "},
		{lexer.TokenString, "array"},
		{lexer.TokenColon, ":"},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenSquareBracketOpen, "["},
		{lexer.TokenNumber, "1"},
		{lexer.TokenComma, ","},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenNumber, "2"},
		{lexer.TokenComma, ","},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenNumber, "3"},
		{lexer.TokenSquareBracketClose, "]"},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "        "},
		{lexer.TokenCurlyBraceClose, "}"},
		{lexer.TokenComma, ","},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "        "},
		{lexer.TokenString, "comment"},
		{lexer.TokenColon, ":"},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenComment, "/* multi-line\n                      comment */"},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "                    "},
		{lexer.TokenComment, "// single-line comment"},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "                    "},
		{lexer.TokenString, "end"},
		{lexer.TokenColon, ":"},
		{lexer.TokenWhitespace, " "},
		{lexer.TokenInfinity, "Infinity"},
		{lexer.TokenNewLine, "\n"},
		{lexer.TokenWhitespace, "    "}, // Changed to 4 spaces to match actual input
		{lexer.TokenCurlyBraceClose, "}"},
		{lexer.TokenEOF, ""},
	}

	l := lexer.NewLexer(input)

	for i, tt := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("Token %d - Expected type %s, got %s", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("Token %d - Expected literal %q, got %q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
