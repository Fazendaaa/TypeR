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
