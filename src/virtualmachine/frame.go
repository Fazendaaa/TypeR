package virtualmachine

import (
	"../code"
	"../object"
)

// Frame :
type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

// Instructions :
func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}

// InitializeFrame :
func InitializeFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}
}
