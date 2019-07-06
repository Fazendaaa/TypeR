package evaluator

import (
	"fmt"

	"../object"
)

// wrongArgumentsError :
func wrongArgumentsError(size int) object.Object {
	return newError("wrong number of parameters, got=%d, expected=1", size)
}

var builtins = map[string]*object.Builtin{
	"len": object.GetBuiltinByName("len"),
	"first": &object.Builtin{
		Fn: func(arguments ...object.Object) object.Object {
			if 1 != len(arguments) {
				return wrongArgumentsError(len(arguments))
			}

			if object.ARRAY_OBJECT != arguments[0].Type() {
				return newError("argument to `first` must be ARRAY, got %s", arguments[0].Type())
			}

			array := arguments[0].(*object.Array)

			if 0 < len(array.Elements) {
				return array.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(arguments ...object.Object) object.Object {
			if 1 != len(arguments) {
				return wrongArgumentsError(len(arguments))
			}

			if object.ARRAY_OBJECT != arguments[0].Type() {
				return newError("argument to `last` must be ARRAY, got %s", arguments[0].Type())
			}

			array := arguments[0].(*object.Array)
			length := len(array.Elements)

			if 0 < length {
				return array.Elements[length-1]
			}

			return NULL
		},
	},
	"tail": &object.Builtin{
		Fn: func(arguments ...object.Object) object.Object {
			if 1 != len(arguments) {
				return wrongArgumentsError(len(arguments))
			}

			if object.ARRAY_OBJECT != arguments[0].Type() {
				return newError("argument to `tail` must be ARRAY, got %s", arguments[0].Type())
			}

			array := arguments[0].(*object.Array)
			length := len(array.Elements)

			if 0 < length {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, array.Elements[1:length])

				return &object.Array{
					Elements: newElements,
				}
			}

			return NULL
		},
	},
	"push": &object.Builtin{
		Fn: func(arguments ...object.Object) object.Object {
			if 2 != len(arguments) {
				return wrongArgumentsError(len(arguments))
			}

			if object.ARRAY_OBJECT != arguments[0].Type() {
				return newError("argument to `push` must be ARRAY, got %s", arguments[0].Type())
			}

			array := arguments[0].(*object.Array)
			length := len(array.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, array.Elements)
			newElements[length] = arguments[1]

			return &object.Array{
				Elements: newElements,
			}
		},
	},
	"puts": &object.Builtin{
		Fn: func(arguments ...object.Object) object.Object {
			for _, argument := range arguments {
				fmt.Println(argument.Inspect())
			}

			return NULL
		},
	},
}
