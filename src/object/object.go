package object

import (
	"fmt"
)

// ObjectType :
type ObjectType string

const (
	INTEGER_OBJECT      = "INTEGER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NULL_OBJECT         = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
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
