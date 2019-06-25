package evaluator

import "../object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(arguments ...object.Object) object.Object {
			if 1 != len(arguments) {
				return newError("wrong number of arguments, got=%d, expected=1", len(arguments))
			}

			switch argument := arguments[0].(type) {
			case *object.String:
				return &object.Integer{
					Value: int64(len(argument.Value)),
				}
			default:
				return newError("argument to `len` not supported, got %s", arguments[0].Type())
			}
		},
	},
}
