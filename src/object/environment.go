package object

// InitializeEnvironment :
func InitializeEnvironment() *Environment {
	store := make(map[string]Object)

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
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && nil != e.outer {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set :
func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value

	return value
}
