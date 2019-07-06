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
