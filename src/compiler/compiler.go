package compiler

import (
	"../ast"
	"../code"
	"../object"
)

// Compiler :
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

// Bytecode :
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// InitializeCompiler :
func InitializeCompiler() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// addInstruction :
func (c *Compiler) addInstruction(instructions []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, instructions...)

	return posNewInstruction
}

// emit :
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	instructions := code.Make(op, operands...)
	position := c.addInstruction(instructions)

	return position
}

// addConstant :
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)

	return len(c.constants) - 1
}

// Compile :
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, statement := range node.Statements {
			err := c.Compile(statement)

			if nil != err {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)

		if nil != err {
			return err
		}

	case *ast.InfixExpression:
		err := c.Compile(node.Left)

		if nil != err {
			return err
		}

		err = c.Compile(node.Right)

		if nil != err {
			return err
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{
			Value: node.Value,
		}

		c.emit(code.OpConstant, c.addConstant(integer))

		return nil
	}

	return nil
}

// Bytecode :
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
