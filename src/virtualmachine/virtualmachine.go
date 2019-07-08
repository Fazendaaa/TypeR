package virtualmachine

import (
	"fmt"

	"../code"
	"../compiler"
	"../object"
)

// GlobalSize : 2 ^ 16 == 65536
const GlobalSize = 65536

// StackSize :
const StackSize = 2048

// FrameSize :
const FrameSize = 1024

// TRUE :
var TRUE = &object.Boolean{
	Value: true,
}

// FALSE :
var FALSE = &object.Boolean{
	Value: false,
}

// NULL :
var NULL = &object.Null{}

// VirtualMachine :
type VirtualMachine struct {
	constants []object.Object

	stack []object.Object
	sp    int

	globals []object.Object

	frames      []*Frame
	framesIndex int
}

// nativeBoolToBooleanObject :
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

// isTruthy :
func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

// currentFrame :
func (vm *VirtualMachine) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

// pushFrame :
func (vm *VirtualMachine) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

// popFrame :
func (vm *VirtualMachine) popFrame() *Frame {
	vm.framesIndex--

	return vm.frames[vm.framesIndex]
}

// LastPoppedStackElement :
func (vm *VirtualMachine) LastPoppedStackElement() object.Object {
	return vm.stack[vm.sp]
}

// StackTop :
func (vm *VirtualMachine) StackTop() object.Object {
	if 0 == vm.sp {
		return nil
	}

	return vm.stack[vm.sp-1]
}

// push :
func (vm *VirtualMachine) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp++

	return nil
}

// pop :
func (vm *VirtualMachine) pop() object.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp--

	return obj
}

// executeIntegerBinaryOperation :
func (vm *VirtualMachine) executeIntegerBinaryOperation(op code.Opcode, left, right object.Object) error {
	var result int64

	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSubtract:
		result = leftValue - rightValue
	case code.OpMultiply:
		result = leftValue * rightValue
	case code.OpDivide:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{
		Value: result,
	})
}

// executeStringBinaryOperation :
func (vm *VirtualMachine) executeStringBinaryOperation(op code.Opcode, left, right object.Object) error {
	if code.OpAdd != op {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{
		Value: leftValue + rightValue,
	})
}

// executeBinaryOperation :
func (vm *VirtualMachine) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {
	case object.INTEGER_OBJECT == leftType && object.INTEGER_OBJECT == rightType:
		return vm.executeIntegerBinaryOperation(op, left, right)
	case object.STRING_OBJECT == leftType && object.STRING_OBJECT == rightType:
		return vm.executeStringBinaryOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
	}
}

// executeIntegerComparisson :
func (vm *VirtualMachine) executeIntegerComparisson(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

// executeComparison :
func (vm *VirtualMachine) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if object.INTEGER_OBJECT == left.Type() && object.INTEGER_OBJECT == right.Type() {
		return vm.executeIntegerComparisson(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d %s %s", op, left.Type(), right.Type())
	}
}

// executeBangOperator :
func (vm *VirtualMachine) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case TRUE:
		return vm.push(FALSE)
	case FALSE:
		return vm.push(TRUE)
	case NULL:
		return vm.push(TRUE)
	default:
		return vm.push(FALSE)
	}
}

// executeMinusOperator :
func (vm *VirtualMachine) executeMinusOperator() error {
	operand := vm.pop()

	if object.INTEGER_OBJECT != operand.Type() {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.Integer).Value

	return vm.push(&object.Integer{
		Value: -value,
	})
}

// buildArray :
func (vm *VirtualMachine) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for index := startIndex; index < endIndex; index++ {
		elements[index-startIndex] = vm.stack[index]
	}

	return &object.Array{
		Elements: elements,
	}
}

// executeArrayIndex :
func (vm *VirtualMachine) executeArrayIndex(array, index object.Object) error {
	arrayObject := array.(*object.Array)
	position := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if position < 0 || position > max {
		return vm.push(NULL)
	}

	return vm.push(arrayObject.Elements[position])
}

// executeIndexExpression :
func (vm *VirtualMachine) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJECT && index.Type() == object.INTEGER_OBJECT:
		return vm.executeArrayIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

// callClosure :
func (vm *VirtualMachine) callClosure(cl *object.Closure, numberOfParameters int) error {
	if numberOfParameters != cl.Fn.NumberOfParameters {
		return fmt.Errorf("wrong number of parameters: want=%d, got=%d", cl.Fn.NumberOfParameters, numberOfParameters)
	}

	frame := InitializeFrame(cl, vm.sp-numberOfParameters)
	vm.pushFrame(frame)

	vm.sp = frame.basePointer + cl.Fn.NumberOfLocals

	return nil
}

// callBuiltin :
func (vm *VirtualMachine) callBuiltin(builtin *object.Builtin, numberOfParameters int) error {
	parameters := vm.stack[vm.sp-numberOfParameters : vm.sp]
	result := builtin.Fn(parameters...)
	vm.sp = vm.sp - numberOfParameters - 1

	if nil != result {
		vm.push(result)
	} else {
		vm.push(NULL)
	}

	return nil
}

// exectueCall :
func (vm *VirtualMachine) exectueCall(numberOfParameters int) error {
	callee := vm.stack[vm.sp-1-numberOfParameters]

	// fmt.Println(callee)

	switch calleeType := callee.(type) {
	case *object.Closure:
		return vm.callClosure(calleeType, numberOfParameters)
	case *object.Builtin:
		return vm.callBuiltin(calleeType, numberOfParameters)
	default:
		return fmt.Errorf("calling a non-function and non-built-in")
	}
}

// pushClosure :
func (vm *VirtualMachine) pushClosure(constIndex int, numberOfFreeVariables int) error {
	constant := vm.constants[constIndex]
	function, ok := constant.(*object.CompiledFunction)

	if !ok {
		return fmt.Errorf("not a function: %+v", constant)
	}

	freeVariables := make([]object.Object, numberOfFreeVariables)

	for index := 0; index < numberOfFreeVariables; index++ {
		freeVariables[index] = vm.stack[vm.sp-numberOfFreeVariables+index]
	}

	vm.sp = vm.sp - numberOfFreeVariables

	closure := &object.Closure{
		Fn:            function,
		FreeVariables: freeVariables,
	}

	return vm.push(closure)
}

// Run :
func (vm *VirtualMachine) Run() error {
	var ip int
	var instructions code.Instructions
	var op code.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		instructions = vm.currentFrame().Instructions()
		op = code.Opcode(instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(instructions[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.constants[constIndex])

			if nil != err {
				return err
			}

		case code.OpAdd, code.OpSubtract, code.OpMultiply, code.OpDivide:
			err := vm.executeBinaryOperation(op)

			if nil != err {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpTrue:
			err := vm.push(TRUE)

			if nil != err {
				return err
			}

		case code.OpFalse:
			err := vm.push(FALSE)

			if nil != err {
				return err
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)

			if nil != err {
				return err
			}

		case code.OpBang:
			err := vm.executeBangOperator()

			if nil != err {
				return err
			}

		case code.OpMinus:
			err := vm.executeMinusOperator()

			if nil != err {
				return err
			}

		case code.OpJump:
			position := int(code.ReadUint16(instructions[ip+1:]))
			vm.currentFrame().ip = position - 1

		case code.OpJumpNotTruthy:
			position := int(code.ReadUint16(instructions[ip+1:]))
			vm.currentFrame().ip += 2

			condition := vm.pop()

			if !isTruthy(condition) {
				vm.currentFrame().ip = position - 1
			}

		case code.OpNull:
			err := vm.push(NULL)

			if nil != err {
				return err
			}

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(instructions[ip+1:])
			vm.currentFrame().ip += 2

			vm.globals[globalIndex] = vm.pop()

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(instructions[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])

			if nil != err {
				return err
			}

		case code.OpArray:
			numberElements := int(code.ReadUint16(instructions[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArray(vm.sp-numberElements, vm.sp)
			vm.sp = vm.sp - numberElements

			err := vm.push(array)

			if nil != err {
				return err
			}

		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeIndexExpression(left, index)

			if nil != err {
				return err
			}

		case code.OpCall:
			numberOfParameters := code.ReadUint8(instructions[ip+1:])

			vm.currentFrame().ip++

			err := vm.exectueCall(int(numberOfParameters))

			if nil != err {
				return err
			}

		case code.OpReturnValue:
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)

			if nil != err {
				return err
			}

		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(NULL)

			if nil != err {
				return err
			}

		case code.OpSetLocal:
			localIndex := code.ReadUint8(instructions[ip+1:])
			vm.currentFrame().ip++

			frame := vm.currentFrame()

			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()

		case code.OpGetLocal:
			localIndex := code.ReadUint8(instructions[ip+1:])
			vm.currentFrame().ip++

			frame := vm.currentFrame()

			err := vm.push(vm.stack[frame.basePointer+int(localIndex)])

			if nil != err {
				return err
			}

		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint8(instructions[ip+1:])

			vm.currentFrame().ip++

			definition := object.Builtins[builtinIndex]

			err := vm.push(definition.Builtin)

			if nil != err {
				return err
			}

		case code.OpClosure:
			constIndex := code.ReadUint16(instructions[ip+1:])
			numberOfFreeVariables := code.ReadUint8(instructions[ip+3:])

			vm.currentFrame().ip += 3

			err := vm.pushClosure(int(constIndex), int(numberOfFreeVariables))

			if nil != err {
				return err
			}

		case code.OpGetFreeVariable:
			freeVariableIndex := code.ReadUint8(instructions[ip+1:])
			vm.currentFrame().ip++

			currentClosure := vm.currentFrame().cl

			err := vm.push(currentClosure.FreeVariables[freeVariableIndex])

			if nil != err {
				return err
			}

		case code.OpCurrentClosure:
			currentClosure := vm.currentFrame().cl
			err := vm.push(currentClosure)

			if nil != err {
				return err
			}

		}
	}

	return nil
}

// InitializeVirtualMachine :
func InitializeVirtualMachine(bytecode *compiler.Bytecode) *VirtualMachine {
	mainFictional := &object.CompiledFunction{
		Instructions: bytecode.Instructions,
	}
	mainClosure := &object.Closure{
		Fn: mainFictional,
	}
	mainFrame := InitializeFrame(mainClosure, 0)

	frames := make([]*Frame, FrameSize)
	frames[0] = mainFrame

	return &VirtualMachine{
		constants: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalSize),

		frames:      frames,
		framesIndex: 1,
	}
}

// InitializeWithGlobalStore :
func InitializeWithGlobalStore(bytecode *compiler.Bytecode, s []object.Object) *VirtualMachine {
	vm := InitializeVirtualMachine(bytecode)
	vm.globals = s

	return vm
}
