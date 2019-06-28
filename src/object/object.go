package object

import (
	"bytes"
	"fmt"
	"strings"

	"../ast"
)

// ObjectType :
type ObjectType string

// BuiltinFunction :
type BuiltinFunction func(arguments ...Object) Object

const (
	INTEGER_OBJECT      = "INTEGER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NULL_OBJECT         = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION"
	STRING_OBJECT       = "STRING"
	BUILTIN_OBJECT      = "BUILTIN"
	ARRAY_OBJECT        = "ARRAY"
)

// Object :
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer :
type Integer struct {
	Value int64
}

// Boolean :
type Boolean struct {
	Value bool
}

// Null :
type Null struct{}

// ReturnValue :
type ReturnValue struct {
	Value Object
}

// Error :
type Error struct {
	Message string
}

// Environment :
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Function :
type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

// String :
type String struct {
	Value string
}

// Builtin :
type Builtin struct {
	Fn BuiltinFunction
}

// Array :
type Array struct {
	Elements []Object
}

// Inspect :
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type :
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJECT
}

// Inspect :
func (b *Boolean) Inspect() string {
	if b.Value {
		return "TRUE"
	}

	return "FALSE"
}

// Type :
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJECT
}

// Inspect :
func (n *Null) Inspect() string {
	return "NULL"
}

// Type :
func (n *Null) Type() ObjectType {
	return NULL_OBJECT
}

// Type :
func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJECT
}

// Inspect :
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Type :
func (e *Error) Type() ObjectType {
	return ERROR_OBJECT
}

// Inspect :
func (e *Error) Inspect() string {
	return "[ERROR]: " + e.Message
}

// Type :
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJECT
}

// Inspect :
func (f *Function) Inspect() string {
	var out bytes.Buffer
	parameters := []string{}

	for _, parameter := range f.Parameters {
		parameters = append(parameters, parameter.String())
	}

	out.WriteString("function")
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// Type :
func (s *String) Type() ObjectType {
	return STRING_OBJECT
}

// Inspect :
func (s *String) Inspect() string {
	return s.Value
}

// Type :
func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJECT
}

// Inspect :
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// Type :
func (a *Array) Type() ObjectType {
	return ARRAY_OBJECT
}

// Inspect :
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}

	for _, element := range elements {
		elements = append(elements, element)
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
