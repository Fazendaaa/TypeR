package object

// Field :
type Field struct {
	Constant bool
	Value    Object
}

// Environment :
type Environment struct {
	store       map[string]Field
	outer       *Environment
	memoization map[string]Object
}

// InitializeEnvironment :
func InitializeEnvironment() *Environment {
	store := make(map[string]Field)
	memoization := make(map[string]Object)

	return &Environment{
		store:       store,
		outer:       nil,
		memoization: memoization,
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

// GetMemoization :
func (e *Environment) GetMemoization(name string) (Object, bool) {
	obj, ok := e.memoization[name]

	if !ok && nil != e.outer {
		obj, ok = e.outer.GetMemoization(name)
	}

	return obj, ok
}

// SetMemoization :
func (e *Environment) SetMemoization(name string, obj Object) Object {
	e.memoization[name] = obj

	if nil != e.outer {
		e.outer.SetMemoization(name, obj)
	}

	return obj
}
