package compiler

import (
	"fmt"

	"../ast"
	"../code"
	"../object"
)

// Compiler :
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction

	symbolTable *SymbolTable
}

// Bytecode :
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// EmittedInstruction :
type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

// addInstruction :
func (c *Compiler) addInstruction(instructions []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, instructions...)

	return posNewInstruction
}

// setLastInstruction :
func (c *Compiler) setLastInstruction(op code.Opcode, position int) {
	previous := c.lastInstruction
	last := EmittedInstruction{
		Opcode:   op,
		Position: position,
	}

	c.previousInstruction = previous
	c.lastInstruction = last
}

// emit :
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	instructions := code.Make(op, operands...)
	position := c.addInstruction(instructions)

	c.setLastInstruction(op, position)

	return position
}

// addConstant :
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)

	return len(c.constants) - 1
}

// lastInstructionIsPop :
func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

// removeLastPop :
func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

// replaceInstruction :
func (c *Compiler) replaceInstruction(position int, newInstruction []byte) {
	for index := 0; index < len(newInstruction); index++ {
		c.instructions[position+index] = newInstruction[index]
	}
}

// changeOperand :
func (c *Compiler) changeOperand(opPosition int, operand int) {
	op := code.Opcode(c.instructions[opPosition])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPosition, newInstruction)
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

		c.emit(code.OpPop)

	case *ast.InfixExpression:
		if "<" == node.Operator {
			err := c.Compile(node.Right)

			if nil != err {
				return err
			}

			err = c.Compile(node.Left)

			if nil != err {
				return err
			}

			c.emit(code.OpGreaterThan)

			return nil
		}

		err := c.Compile(node.Left)

		if nil != err {
			return err
		}

		err = c.Compile(node.Right)

		if nil != err {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSubtract)
		case "*":
			c.emit(code.OpMultiply)
		case "/":
			c.emit(code.OpDivide)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{
			Value: node.Value,
		}

		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)

		if nil != err {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.ConditionalExpression:
		err := c.Compile(node.Condition)

		if nil != err {
			return err
		}

		jumpNotTruthyPosition := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)

		if nil != err {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		jumpPosition := c.emit(code.OpJump, 9999)

		afterConsequencePosition := len(c.instructions)
		c.changeOperand(jumpNotTruthyPosition, afterConsequencePosition)

		if nil == node.Alternative {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)

			if nil != err {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}

		}

		afterAlternativePosition := len(c.instructions)
		c.changeOperand(jumpPosition, afterAlternativePosition)

	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			err := c.Compile(statement)

			if nil != err {
				return err
			}
		}

	case *ast.LetStatement:
		err := c.Compile(node.Value)

		if nil != err {
			return err
		}

		symbol := c.symbolTable.Define(node.Name.Value)

		c.emit(code.OpSetGlobal, symbol.Index)

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)

		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.emit(code.OpGetGlobal, symbol.Index)

	case *ast.StringLiteral:
		str := &object.String{
			Value: node.Value,
		}

		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, element := range node.Elements {
			err := c.Compile(element)

			if nil != err {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))
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

// InitializeCompiler :
func InitializeCompiler() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         InitializeSymbolTable(),
	}
}

// InitializeWithState :
func InitializeWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := InitializeCompiler()
	compiler.symbolTable = s
	compiler.constants = constants

	return compiler
}
