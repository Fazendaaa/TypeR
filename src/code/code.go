package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions :
type Instructions []byte

// Opcode :
type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSubtract
	OpMultiply
	OpDivide
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
	OpNull
	OpGetGlobal
	OpSetGlobal
	OpArray
	OpIndex
)

// Definition :
type Definition struct {
	Name          string
	OperandWidths []int
}

// definitions :
var definitions = map[Opcode]*Definition{
	OpConstant: {
		"OpConstant",
		[]int{2},
	},
	OpAdd: {
		"OpAdd",
		[]int{},
	},
	OpPop: {
		"OpPop",
		[]int{},
	},
	OpSubtract: {
		"OpSubtract",
		[]int{},
	},
	OpMultiply: {
		"OpMultiply",
		[]int{},
	},
	OpDivide: {
		"OpDivide",
		[]int{},
	},
	OpTrue: {
		"OpTrue",
		[]int{},
	},
	OpFalse: {
		"OpFalse",
		[]int{},
	},
	OpEqual: {
		"OpEqual",
		[]int{},
	},
	OpNotEqual: {
		"OpNotEqual",
		[]int{},
	},
	OpGreaterThan: {
		"OpGreaterThan",
		[]int{},
	},
	OpMinus: {
		"OpMinus",
		[]int{},
	},
	OpBang: {
		"OpBang",
		[]int{},
	},
	OpJumpNotTruthy: {
		"OpJumpNotTruthy",
		[]int{
			2,
		},
	},
	OpJump: {
		"OpJump",
		[]int{
			2,
		},
	},
	OpNull: {
		"OpNull",
		[]int{},
	},
	OpGetGlobal: {
		"OpGetGlobal",
		[]int{
			2,
		},
	},
	OpSetGlobal: {
		"OpSetGlobal",
		[]int{
			2,
		},
	},
	OpArray: {
		"OpArray",
		[]int{
			2,
		},
	},
	OpIndex: {
		"OpIndex",
		[]int{},
	},
}

// Lookup :
func Lookup(op byte) (*Definition, error) {
	definition, ok := definitions[Opcode(op)]

	if !ok {
		return nil, fmt.Errorf("opcode %d is undefined", op)
	}

	return definition, nil
}

// Make :
func Make(op Opcode, operands ...int) []byte {
	definition, ok := definitions[op]

	if !ok {
		return []byte{}
	}

	instructionsLen := 1

	for _, width := range definition.OperandWidths {
		instructionsLen += width
	}

	instruction := make([]byte, instructionsLen)
	instruction[0] = byte(op)
	offset := 1

	for index, operand := range operands {
		width := definition.OperandWidths[index]

		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}

		offset += width
	}

	return instruction
}

// ReadUint16 :
func ReadUint16(instructions Instructions) uint16 {
	return binary.BigEndian.Uint16(instructions)
}

// ReadOperands :
func ReadOperands(definition *Definition, instructions Instructions) ([]int, int) {
	operands := make([]int, len(definition.OperandWidths))
	offset := 0

	for index, width := range definition.OperandWidths {
		switch width {
		case 2:
			operands[index] = int(ReadUint16(instructions[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// fmtInstruction :
func (i Instructions) fmtInstruction(definition *Definition, operands []int) string {
	operandCount := len(definition.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return definition.Name
	case 1:
		return fmt.Sprintf("%s %d", definition.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", definition.Name)
}

// String :
func (i Instructions) String() string {
	var out bytes.Buffer
	index := 0

	for index < len(i) {
		definition, err := Lookup(i[index])

		if nil != err {
			fmt.Fprintf(&out, "ERROR: %s\n", err)

			continue
		}

		operands, read := ReadOperands(definition, i[index+1:])

		fmt.Fprintf(&out, "%04d %s\n", index, i.fmtInstruction(definition, operands))

		index += 1 + read
	}

	return out.String()
}
