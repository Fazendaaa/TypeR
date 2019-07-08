package compiler

// SymbolScope :
type SymbolScope string

const (
	GlobalScope       SymbolScope = "GLOBAL"
	LocalScope        SymbolScope = "LOCAL"
	BuiltinScope      SymbolScope = "BUILTIN"
	FreeVariableScope SymbolScope = "FREE_VARIABLE"
	FunctionScope     SymbolScope = "FUNCTION"
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

	FreeVariableSymbol []Symbol
}

// defineFreeVariable :
func (s *SymbolTable) defineFreeVariable(original Symbol) Symbol {
	s.FreeVariableSymbol = append(s.FreeVariableSymbol, original)

	symbol := Symbol{
		Name:  original.Name,
		Index: len(s.FreeVariableSymbol) - 1,
	}
	symbol.Scope = FreeVariableScope

	s.store[original.Name] = symbol

	return symbol
}

// DefineFunctionName :
func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Index: 0,
		Scope: FunctionScope,
	}

	s.store[name] = symbol

	return symbol
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

		if !ok {
			return obj, ok
		}

		if GlobalScope == obj.Scope || BuiltinScope == obj.Scope {
			return obj, ok
		}

		freeVariable := s.defineFreeVariable(obj)

		return freeVariable, true
	}

	return obj, ok
}

// InitializeSymbolTable :
func InitializeSymbolTable() *SymbolTable {
	store := make(map[string]Symbol)
	freeVariable := []Symbol{}

	return &SymbolTable{
		store:              store,
		FreeVariableSymbol: freeVariable,
	}
}

// InitializeEnclosedSymbolTable :
func InitializeEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := InitializeSymbolTable()

	s.Outer = outer

	return s
}
