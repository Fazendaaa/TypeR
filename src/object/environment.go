package object

// Field :
type Field struct {
	Constant bool
	Value    Object
}

// Environment :
type Environment struct {
	store map[string]Field
	outer *Environment
}

// InitializeEnvironment :
func InitializeEnvironment() *Environment {
	store := make(map[string]Field)

	return &Environment{
		store: store,
		outer: nil,
	}
}

// InitializeEnclosedEnvironment :
func InitializeEnclosedEnvironment(outer *Environment) *Environment {
	environment := InitializeEnvironment()
	environment.outer = outer

	return environment
}

// Get :
func (e *Environment) Get(name string) (Field, bool) {
	obj, ok := e.store[name]

	if !ok && nil != e.outer {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set :
func (e *Environment) Set(name string, constant bool, value Object) Object {
	field := Field{
		Value:    value,
		Constant: constant,
	}

	e.store[name] = field

	return value
}
