package evaluator

import (
	"fmt"

	"../ast"
	"../object"
)

var (
	NULL = &object.Null{}
	TRUE = &object.Boolean{
		Value: true,
	}
	FALSE = &object.Boolean{
		Value: false,
	}
)

// isError :
func isError(obj object.Object) bool {
	if nil != obj {
		return object.ERROR_OBJECT == obj.Type()
	}

	return false
}

// newError :
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

// evalProgram :
func evalProgram(statements []ast.Statement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, environment)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

// nativeBoolToBooleanObject :
func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}

	return FALSE
}

// evalBangOperatorExpression :
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

// evalMinusPrefixOperatorExpression :
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJECT {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value

	return &object.Integer{
		Value: -value,
	}
}

// evalPrefixExpression :
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// evalIntgerInfixExpression :
func evalIntgerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{
			Value: leftValue + rightValue,
		}
	case "-":
		return &object.Integer{
			Value: leftValue - rightValue,
		}
	case "*":
		return &object.Integer{
			Value: leftValue * rightValue,
		}
	case "/":
		return &object.Integer{
			Value: leftValue / rightValue,
		}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalInfixExpression :
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalIntgerInfixExpression(operator, left, right)
	case "==" == operator:
		return nativeBoolToBooleanObject(left == right)
	case "!=" == operator:
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// isTruthy :
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

// evalConditionalExpression :
func evalConditionalExpression(ce *ast.ConditionalExpression, environment *object.Environment) object.Object {
	condition := Eval(ce.Condition, environment)

	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ce.Consequence, environment)
	} else if nil != ce.Alternative {
		return Eval(ce.Alternative, environment)
	}

	return NULL
}

// evalBlockStatement :
func evalBlockStatement(block *ast.BlockStatement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, environment)

		if nil != result {
			resultType := result.Type()

			if resultType == object.RETURN_VALUE_OBJECT || resultType == object.ERROR_OBJECT {
				return result
			}
		}
	}

	return result
}

// evalIdentifier :
func evalIdentifier(node *ast.Identifier, environment *object.Environment) object.Object {
	value, ok := environment.Get(node.Value)

	if !ok {
		return newError("identifier not found: " + node.Value)
	}

	return value
}

// evalExpression :
func evalExpression(expressions []ast.Expression, environment *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := Eval(expression, environment)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

// unwrapReturnValue :
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

// extendendFunctionEnvironment :
func extendendFunctionEnvironment(fn *object.Function, arguments []object.Object) *object.Environment {
	environment := object.InitializeEnclosedEnvironment(fn.Environment)

	for parameterIdx, parameter := range fn.Parameters {
		environment.Set(parameter.Value, arguments[parameterIdx])
	}

	return environment
}

// applyFunction :
func applyFunction(fn object.Object, arguments []object.Object) object.Object {
	function, ok := fn.(*object.Function)

	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnvironment := extendendFunctionEnvironment(function, arguments)
	evaluated := Eval(function.Body, extendedEnvironment)

	return unwrapReturnValue(evaluated)
}

// Eval :
func Eval(node ast.Node, environment *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, environment)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, environment)
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, environment)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, environment)

		if isError(left) {
			return left
		}

		right := Eval(node.Right, environment)

		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, environment)
	case *ast.ConditionalExpression:
		return evalConditionalExpression(node, environment)
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, environment)

		if isError(value) {
			return value
		}

		return &object.ReturnValue{
			Value: value,
		}
	case *ast.LetStatement:
		value := Eval(node.Value, environment)

		if isError(value) {
			return value
		}

		environment.Set(node.Name.Value, value)
	case *ast.Identifier:
		return evalIdentifier(node, environment)
	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body

		return &object.Function{
			Parameters:  parameters,
			Environment: environment,
			Body:        body,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, environment)

		if isError(function) {
			return function
		}

		arguments := evalExpression(node.Arguments, environment)

		if 1 == len(arguments) && isError(arguments[0]) {
			return arguments[0]
		}

		return applyFunction(function, arguments)
	}

	return nil
}
