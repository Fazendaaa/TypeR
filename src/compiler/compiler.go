package compiler

import (
	"fmt"

	"../ast"
	"../code"
	"../object"
)

// CompilationsScope :
type CompilationsScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

// Compiler :
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object

	symbolTable *SymbolTable

	scopes     []CompilationsScope
	scopeIndex int
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

// enterScope :
func (c *Compiler) enterScope() {
	scope := CompilationsScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = InitializeEnclosedSymbolTable(c.symbolTable)
}

// leaveScope :
func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

// currentInstructions :
func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

// addInstruction :
func (c *Compiler) addInstruction(instructions []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), instructions...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

// setLastInstruction :
func (c *Compiler) setLastInstruction(op code.Opcode, position int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{
		Opcode:   op,
		Position: position,
	}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
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

// lastInstructionIs :
func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

// removeLastPop :
func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

// replaceInstruction :
func (c *Compiler) replaceInstruction(position int, newInstruction []byte) {
	instruction := c.currentInstructions()

	for index := 0; index < len(newInstruction); index++ {
		instruction[position+index] = newInstruction[index]
	}
}

// replaceLastPopWithReturn :
func (c *Compiler) replaceLastPopWithReturn() {
	lastPosition := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPosition, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

// changeOperand :
func (c *Compiler) changeOperand(opPosition int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPosition])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPosition, newInstruction)
}

// loadSymbol :
func (c *Compiler) loadSymbol(symbol Symbol) {
	switch symbol.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, symbol.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, symbol.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, symbol.Index)
	case FreeVariableScope:
		c.emit(code.OpGetFreeVariable, symbol.Index)
	case FunctionScope:
		c.emit(code.OpCurrentClosure)
	}
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

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		jumpPosition := c.emit(code.OpJump, 9999)

		afterConsequencePosition := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPosition, afterConsequencePosition)

		if nil == node.Alternative {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)

			if nil != err {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}

		}

		afterAlternativePosition := len(c.currentInstructions())
		c.changeOperand(jumpPosition, afterAlternativePosition)

	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			err := c.Compile(statement)

			if nil != err {
				return err
			}
		}

	case *ast.LetStatement:
		symbol := c.symbolTable.Define(node.Name.Value)
		err := c.Compile(node.Value)

		if nil != err {
			return err
		}

		if GlobalScope == symbol.Scope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)

		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.loadSymbol(symbol)

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

	case *ast.IndexExpression:
		err := c.Compile(node.Left)

		if nil != err {
			return err
		}

		err = c.Compile(node.Index)

		if nil != err {
			return err
		}

		c.emit(code.OpIndex)

	case *ast.FunctionLiteral:
		c.enterScope()

		if "" != node.Name {
			c.symbolTable.DefineFunctionName(node.Name)
		}

		for _, parameter := range node.Parameters {
			c.symbolTable.Define(parameter.Value)
		}

		err := c.Compile(node.Body)

		if nil != err {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}

		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		freeVariableSymbols := c.symbolTable.FreeVariableSymbol
		numberOfLocals := c.symbolTable.numberDefinitions
		instructions := c.leaveScope()

		for _, symbol := range freeVariableSymbols {
			c.loadSymbol(symbol)
		}

		compiledFunction := &object.CompiledFunction{
			Instructions:       instructions,
			NumberOfLocals:     numberOfLocals,
			NumberOfParameters: len(node.Parameters),
		}

		functionIndex := c.addConstant(compiledFunction)
		c.emit(code.OpClosure, functionIndex, len(freeVariableSymbols))

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)

		if nil != err {
			return err
		}

		c.emit(code.OpReturnValue)

	case *ast.CallExpression:
		err := c.Compile(node.Function)

		if nil != err {
			return err
		}

		for _, parameter := range node.Parameters {
			err := c.Compile(parameter)

			if nil != err {
				return err
			}
		}

		c.emit(code.OpCall, len(node.Parameters))

	}

	return nil
}

// Bytecode :
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

// InitializeCompiler :
func InitializeCompiler() *Compiler {
	mainScope := CompilationsScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := InitializeSymbolTable()

	for index, value := range object.Builtins {
		symbolTable.DefineBuiltin(index, value.Name)
	}

	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
		symbolTable:  symbolTable,
		scopes:       []CompilationsScope{mainScope},
		scopeIndex:   0,
	}
}

// InitializeWithState :
func InitializeWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := InitializeCompiler()
	compiler.symbolTable = s
	compiler.constants = constants

	return compiler
}
