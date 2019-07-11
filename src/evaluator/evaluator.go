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

// evalStringInfixExpression :
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if "+" != operator {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return &object.String{
		Value: leftValue + rightValue,
	}
}

// evalInfixExpression :
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalIntgerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJECT && right.Type() == object.STRING_OBJECT:
		return evalStringInfixExpression(operator, left, right)
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

// evalConstant :
func evalConstant(cons *ast.ConstStatement, environment *object.Environment) object.Object {
	if _, ok := environment.Get(cons.Name.Value); ok {
		return newError("constant '%s' value cannot be overwritten", cons.Name.Value)
	}

	if _, ok := builtins[cons.Name.Value]; ok {
		return newError("builtin function '%s' cannot be overwritten", cons.Name.Value)
	}

	value := Eval(cons.Value, environment)

	if isError(value) {
		return value
	}

	environment.Set(cons.Name.Value, value)

	return nil
}

// evalIdentifier :
func evalIdentifier(node *ast.Identifier, environment *object.Environment) object.Object {
	if value, ok := environment.Get(node.Value); ok {
		return value
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
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

// evalArrayIndexExpression :
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	position := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if position < 0 || position > max {
		return NULL
	}

	return arrayObject.Elements[position]
}

// evalIndexExpression :
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJECT && index.Type() == object.INTEGER_OBJECT:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
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
	switch function := fn.(type) {
	case *object.Function:
		extendendEnvironment := extendendFunctionEnvironment(function, arguments)
		evaluated := Eval(function.Body, extendendEnvironment)

		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		if result := function.Fn(arguments...); nil != result {
			return result
		}

		return NULL
	default:
		return newError("not a function: %s", fn.Type())
	}
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

	case *ast.ConstStatement:
		err := evalConstant(node, environment)

		if nil != err {
			return err
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

		arguments := evalExpression(node.Parameters, environment)

		if 1 == len(arguments) && isError(arguments[0]) {
			return arguments[0]
		}

		return applyFunction(function, arguments)

	case *ast.StringLiteral:
		return &object.String{
			Value: node.Value,
		}

	case *ast.ArrayLiteral:
		elements := evalExpression(node.Elements, environment)

		if 1 == len(elements) && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{
			Elements: elements,
		}

	case *ast.IndexExpression:
		left := Eval(node.Left, environment)

		if isError(left) {
			return left
		}

		index := Eval(node.Index, environment)

		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	}

	return nil
}
