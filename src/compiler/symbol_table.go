package compiler

// SymbolScope :
type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
)

// Symbol :
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable :
type SymbolTable struct {
	Outer *SymbolTable

	store             map[string]Symbol
	numberDefinitions int
}

// Define :
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Index: s.numberDefinitions,
	}

	if nil == s.Outer {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numberDefinitions++

	return symbol
}

// DefineBuiltin :
func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Index: index,
		Scope: BuiltinScope,
	}

	s.store[name] = symbol

	return symbol
}

// Resolve :
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	if !ok && nil != s.Outer {
		obj, ok = s.Outer.Resolve(name)
	}

	return obj, ok
}

// InitializeSymbolTable :
func InitializeSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)

	return &SymbolTable{
		store: s,
	}
}

// InitializeEnclosedSymbolTable :
func InitializeEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := InitializeSymbolTable()

	s.Outer = outer

	return s
}
