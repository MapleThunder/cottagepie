package evaluator

import (
	"cottagepie/lexer"
	"cottagepie/object"
	"cottagepie/parser"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"Unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
				if (10 > 1) {
					if (10 > 1) {
						serves true + false;
					}

					serves 1;
				}
			`,
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"Identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("No error object served. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("Wrong error message. Expected=%s, got=%s", tt.expectedMessage, errObj.Message)
			continue
		}
	}

}

// Test Expressions
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"7", 7},
		{"22", 22},
		{"-7", -7},
		{"-22", -22},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Cristiano Ronaldo!"`
	expected := "Cristiano Ronaldo!"

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Errorf("Object is not a String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != expected {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// Test Statements
func TestServesStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"serves 10;", 10},
		{"serves 10; 9;", 10},
		{"serves 2 * 5; 9;", 10},
		{"9; serves 2 * 5; 9;", 10},
		{
			`if (10 > 1) {
				if (10 > 1) {
				  serves 10;
				}
				serves 1;
			}`, 10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBakeStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"bake a = 5; a;", 5},
		{"bake a = 5 * 5; a;", 25},
		{"bake a = 5; bake b = a; b;", 5},
		{"bake a = 5; bake b = a; bake c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestRecipeObject(t *testing.T) {
	input := "recipe(x) { x + 7; };"

	evaluated := testEval(input)
	rc, ok := evaluated.(*object.Recipe)
	if !ok {
		t.Fatalf("Object is not Recipe. Got=%T (%+v)", evaluated, evaluated)
	}

	if len(rc.Parameters) != 1 {
		t.Fatalf("Recipe has wrong parameters. Parameters=%+v", rc.Parameters)
	}

	if rc.Parameters[0].String() != "x" {
		t.Fatalf("First parameter is not 'x', got=%q", rc.Parameters[0])
	}

	expectedBody := "(x + 7)"

	if rc.Body.String() != expectedBody {
		t.Fatalf("Body is not %q, got=%q", expectedBody, rc.Body.String())
	}
}

func TestRecipeApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"bake identity to recipe(x) { x; }; identity(5);", 5},
		{"bake identity to recipe(x) { serves x; }; identity(5);", 5},
		{"bake double to recipe(x) { x * 2; }; double(5);", 10},
		{"bake add to recipe(x, y) { x + y; }; add(5, 5);", 10},
		{"bake add to recipe(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"recipe(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		bake newAdder to recipe(x) {
			rc(y) { x + y };
		};

		bake addTwo to newAdder(2);
		addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

// Test Operators
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!7", false},
		{"!!true", true},
		{"!!false", false},
		{"!!7", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	book := object.NewCookbook()

	return Eval(program, book)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
