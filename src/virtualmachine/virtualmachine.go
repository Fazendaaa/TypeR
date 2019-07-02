package virtualmachine

import (
	"fmt"

	"../code"
	"../compiler"
	"../object"
)

// StackSize :
const StackSize = 2048

// VirtualMachine :
type VirtualMachine struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int
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

// executeBinaryOperation :
func (vm *VirtualMachine) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if object.INTEGER_OBJECT == leftType && object.INTEGER_OBJECT == rightType {
		return vm.executeIntegerBinaryOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
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
	}
}
