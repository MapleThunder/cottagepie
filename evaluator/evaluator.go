package evaluator

import (
	"cottagepie/ast"
	"cottagepie/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, book *object.Cookbook) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, book)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, book)

	case *ast.BakeStatement:
		val := Eval(node.Value, book)
		if isError(val) {
			return val
		}
		book.Set(node.Name.Value, val)

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, book)

	case *ast.RecipeLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Recipe{Parameters: params, Body: body, Cookbook: book}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right, book)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, book)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, book)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IndexExpression:
		left := Eval(node.Left, book)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, book)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.IfExpression:
		return evalIfExpression(node, book)

	case *ast.BlockStatement:
		return evalBlockStatement(node, book)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.ServesStatement:
		val := Eval(node.ServesValue, book)
		if isError(val) {
			return val
		}
		return &object.ServesValue{Value: val}

	case *ast.CallExpression:
		recipe := Eval(node.Recipe, book)
		if isError(recipe) {
			return recipe
		}
		args := evalExpressions(node.Arguments, book)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyRecipe(recipe, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, book)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	}

	return nil
}

func evalProgram(program *ast.Program, book *object.Cookbook) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, book)

		switch result := result.(type) {
		case *object.ServesValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("Unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalBlockStatement(block *ast.BlockStatement, book *object.Cookbook) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, book)

		if result != nil {
			rt := result.Type()
			if rt == object.SERVES_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("Unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIfExpression(ie *ast.IfExpression, book *object.Cookbook) object.Object {
	condition := Eval(ie.Condition, book)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, book)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, book)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, book *object.Cookbook) object.Object {
	if val, ok := book.Get(node.Value); ok {
		return val
	}

	if built_in, ok := built_ins[node.Value]; ok {
		return built_in
	}

	return newError("Identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, book *object.Cookbook) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, book)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("Index operator not supported: %s", left.Type())
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func applyRecipe(rc object.Object, args []object.Object) object.Object {
	switch recipe := rc.(type) {
	case *object.Recipe:
		extendedBook := extendRecipeBook(recipe, args)
		evaluated := Eval(recipe.Body, extendedBook)
		return unwrapServesValue(evaluated)

	case *object.BuiltIn:
		return recipe.Fn(args...)

	default:
		return newError("Not a function: %s", rc.Type())
	}
}

func extendRecipeBook(rc *object.Recipe, args []object.Object) *object.Cookbook {
	book := object.NewExtendedCookbook(rc.Cookbook)

	for param_idx, param := range rc.Parameters {
		book.Set(param.Value, args[param_idx])
	}

	return book
}

func unwrapServesValue(obj object.Object) object.Object {
	if serves_value, ok := obj.(*object.ServesValue); ok {
		return serves_value.Value
	}
	return obj
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}
