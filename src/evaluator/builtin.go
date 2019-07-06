package evaluator

import (
	"../object"
)

// wrongArgumentsError :
func wrongArgumentsError(size int) object.Object {
	return newError("wrong number of parameters, got=%d, expected=1", size)
}

var builtins = map[string]*object.Builtin{
	"puts": object.GetBuiltinByName("puts"),
	"len":  object.GetBuiltinByName("len"),
	"head": object.GetBuiltinByName("head"),
	"tail": object.GetBuiltinByName("tail"),
	"last": object.GetBuiltinByName("last"),
	"push": object.GetBuiltinByName("push"),
}
