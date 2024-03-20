package parser

import (
	"cottagepie/ast"
	"cottagepie/lexer"
	"testing"
)

func TestBakeStatements(t *testing.T) {
	input := `bake x to 5;
		bake y to 10;
		bake title to "CottagePie";
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"title"},
	}

	for index, testData := range tests {
		statement := program.Statements[index]
		if !testBakeStatement(t, statement, testData.expectedIdentifier) {
			return
		}
	}
}

func testBakeStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "bake" {
		t.Errorf("s.TokenLiteral not 'bake'. got=%q", statement.TokenLiteral())
		return false
	}

	bakeStatement, ok := statement.(*ast.BakeStatement)
	if !ok {
		t.Errorf("statement not *ast.BakeStatement. got=%T", statement)
		return false
	}

	if bakeStatement.Name.Value != name {
		t.Errorf("bakeStatement.Name.Value not '%s'. got=%s", name, bakeStatement.Name.Value)
	}

	if bakeStatement.Name.TokenLiteral() != name {
		t.Errorf("bakeStatement.Name.TokenLiteral() not '%s'. got=%s", name, bakeStatement.Name.TokenLiteral())
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}
