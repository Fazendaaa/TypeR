package object

import "fmt"

// newError :
func newError(format string, variables ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, variables...),
	}
}

// Builtins :
var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"puts",
		&Builtin{
			Fn: func(parameters ...Object) Object {
				for _, parameter := range parameters {
					fmt.Println(parameter.Inspect())
				}

				return nil
			},
		},
	},
	{
		"len",
		&Builtin{
			Fn: func(parameters ...Object) Object {
				if 1 != len(parameters) {
					return newError("wrong number of parameters, got=%d, want=1", len(parameters))
				}

				switch parameter := parameters[0].(type) {
				case *Array:
					return &Integer{
						Value: int64(len(parameter.Elements)),
					}

				case *String:
					return &Integer{
						Value: int64(len(parameter.Value)),
					}

				default:
					return newError("parameters to `len` not supported, got=%s", parameters[0].Type())
				}
			},
		},
	},
	{
		"head",
		&Builtin{
			Fn: func(parameters ...Object) Object {
				if 1 != len(parameters) {
					return newError("wrong number of parameters, got=%d, want=1", len(parameters))
				}

				if ARRAY_OBJECT != parameters[0].Type() {
					return newError("parameter to `head` must be ARRAY, got %s", parameters[0].Type())
				}

				arr := parameters[0].(*Array)

				if 0 < len(arr.Elements) {
					return arr.Elements[0]
				}

				return nil
			},
		},
	},
	{
		"tail",
		&Builtin{
			Fn: func(parameters ...Object) Object {
				if 1 != len(parameters) {
					return newError("wrong number of parameters, got=%d, want=1", len(parameters))
				}

				if ARRAY_OBJECT != parameters[0].Type() {
					return newError("parameter to `tail` must be ARRAY, got %s", parameters[0].Type())
				}

				arr := parameters[0].(*Array)
				length := len(arr.Elements)

				if 0 < length {
					newElements := make([]Object, length-1, length-1)

					copy(newElements, arr.Elements[1:length])

					return &Array{
						Elements: newElements,
					}
				}

				return nil
			},
		},
	},
	{
		"last",
		&Builtin{
			Fn: func(parameters ...Object) Object {
				if 1 != len(parameters) {
					return newError("wrong number of parameters, got=%d, want=1", len(parameters))
				}

				if ARRAY_OBJECT != parameters[0].Type() {
					return newError("parameter to `last` must be ARRAY, got %s", parameters[0].Type())
				}

				arr := parameters[0].(*Array)
				length := len(arr.Elements)

				if 0 < length {
					return arr.Elements[length-1]
				}

				return nil
			},
		},
	},
	{
		"push",
		&Builtin{
			Fn: func(parameters ...Object) Object {
				if 2 != len(parameters) {
					return newError("wrong number of parameters, got=%d, want=2", len(parameters))
				}

				if ARRAY_OBJECT != parameters[0].Type() {
					return newError("parameter to `push` must be ARRAY, got %s", parameters[0].Type())
				}

				arr := parameters[0].(*Array)
				length := len(arr.Elements)

				newElements := make([]Object, length+1, length+1)

				copy(newElements, arr.Elements[0:length])

				newElements[length] = parameters[1]

				return &Array{
					Elements: newElements,
				}
			},
		},
	},
}

// GetBuiltinByName :
func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}

	return nil
}
