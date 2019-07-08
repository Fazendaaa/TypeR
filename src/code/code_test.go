package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		opcode   Opcode
		operands []int
		expected []byte
	}{
		{
			OpConstant,
			[]int{
				65534,
			},
			[]byte{
				byte(OpConstant),
				255,
				254,
			},
		},
		{
			OpAdd,
			[]int{},
			[]byte{
				byte(OpAdd),
			},
		},
		{
			OpGetLocal,
			[]int{
				255,
			},
			[]byte{
				byte(OpGetLocal),
				255,
			},
		},
		{
			OpClosure,
			[]int{
				65534,
				255,
			},
			[]byte{
				byte(OpClosure),
				255,
				254,
				255,
			},
		},
	}

	for _, tt := range tests {
		instruction := Make(tt.opcode, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length, want=%d, got=%d", len(instruction), len(tt.expected))
		}

		for index, desiredByte := range tt.expected {
			if instruction[index] != tt.expected[index] {
				t.Errorf("wrong byte at pos %d, want=%d, got=%d", index, desiredByte, instruction[index])
			}
		}
	}
}

// TestInstructionsString :
func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpGetLocal, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpClosure, 65535, 255),
	}
	expected := `0000 OpAdd
0001 OpGetLocal 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpClosure 65535 255
`
	concatted := Instructions{}

	for _, instruction := range instructions {
		concatted = append(concatted, instruction...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted\nwant=%q\ngot=%q", expected, concatted.String())
	}

}

// TestReadOperands :
func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{
			OpConstant,
			[]int{
				65535,
			},
			2,
		},
		{
			OpGetLocal,
			[]int{
				255,
			},
			1,
		},
		{
			OpClosure,
			[]int{
				65535,
				255,
			},
			3,
		},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		definition, err := Lookup(byte(tt.op))

		if nil != err {
			t.Fatalf("definition not found: %q", err)
		}

		operandsRead, offset := ReadOperands(definition, instruction[1:])

		if offset != tt.bytesRead {
			t.Fatalf("offset wrong, want=%d, got=%d", tt.bytesRead, offset)
		}

		for index, want := range tt.operands {
			if want != operandsRead[index] {
				t.Errorf("operand wrong, want=%d, got=%d", want, operandsRead[index])
			}
		}
	}
}
