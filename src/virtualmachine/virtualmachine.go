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
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int

	globals []object.Object
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
	postion := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if postion < 0 || postion > max {
		return vm.push(NULL)
	}

	return vm.push(arrayObject.Elements[postion])
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

// Run :
func (vm *VirtualMachine) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

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
			position := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = position - 1

		case code.OpJumpNotTruthy:
			postion := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition := vm.pop()

			if !isTruthy(condition) {
				ip = postion - 1
			}

		case code.OpNull:
			err := vm.push(NULL)

			if nil != err {
				return err
			}

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			vm.globals[globalIndex] = vm.pop()

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.globals[globalIndex])

			if nil != err {
				return err
			}

		case code.OpArray:
			numberElements := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			array := vm.buildArray(vm.sp-numberElements, vm.sp)
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

		}
	}

	return nil
}

// InitializeVirtualMachine :
func InitializeVirtualMachine(bytecode *compiler.Bytecode) *VirtualMachine {
	return &VirtualMachine{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalSize),
	}
}

// InitializeWithGlobalStore :
func InitializeWithGlobalStore(bytecode *compiler.Bytecode, s []object.Object) *VirtualMachine {
	vm := InitializeVirtualMachine(bytecode)
	vm.globals = s

	return vm
}
