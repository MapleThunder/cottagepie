package parser

import (
	"cottagepie/ast"
	"cottagepie/lexer"
	"testing"
)

func TestBakeStatements(testState *testing.T) {
	input := `bake x to 5;
		bake y to 10;
		bake title to "CottagePie";
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		testState.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		testState.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
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
		if !testBakeStatement(testState, statement, testData.expectedIdentifier) {
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
