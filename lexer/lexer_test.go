package lexer

import (
	"testing"

	"cottagepie/token"
)

func TestNextToken(t *testing.T) {
	input := `bake five to 5;
		bake ten to 10;
		
		bake add to recipe(x, y) {
			x + y;
		};
		
		bake result to add(five, ten);
		!-/*5;
		5 < 10 > 5;
		
		if (5 < 10) {
			serves true;
		} else {
			serves false;
		}
		
		10 == 10;
		10 != 9;
		"CR7"
		"Cristiano Ronaldo"
		'CR7'
		'Cristiano Ronaldo'
		[1, 2];
		{"goat": "Cristiano"}
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.BAKE, "bake"},
		{token.IDENT, "five"},
		{token.ASSIGN, "to"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.BAKE, "bake"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "to"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.BAKE, "bake"},
		{token.IDENT, "add"},
		{token.ASSIGN, "to"},
		{token.RECIPE, "recipe"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.BAKE, "bake"},
		{token.IDENT, "result"},
		{token.ASSIGN, "to"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.SERVES, "serves"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.SERVES, "serves"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.STRING, "CR7"},
		{token.STRING, "Cristiano Ronaldo"},
		{token.STRING, "CR7"},
		{token.STRING, "Cristiano Ronaldo"},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.LBRACE, "{"},
		{token.STRING, "goat"},
		{token.COLON, ":"},
		{token.STRING, "Cristiano"},
		{token.RBRACE, "}"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
