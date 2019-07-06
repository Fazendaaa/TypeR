package virtualmachine

import (
	"../code"
	"../object"
)

// Frame :
type Frame struct {
	fn          *object.CompiledFunction
	ip          int
	basePointer int
}

// Instructions :
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}

// InitializeFrame :
func InitializeFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{
		fn:          fn,
		ip:          -1,
		basePointer: basePointer,
	}
}
