package compiler

// SymbolScope :
type SymbolScope string

const (
	// GlobalScope :
	GlobalScope SymbolScope = "GLOBAL"
)

// Symbol :
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable :
type SymbolTable struct {
	store             map[string]Symbol
	numberDefinitions int
}

// InitializeSymbolTable :
func InitializeSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)

	return &SymbolTable{
		store: s,
	}
}

// Define :
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Index: s.numberDefinitions,
		Scope: GlobalScope,
	}

	s.store[name] = symbol
	s.numberDefinitions++

	return symbol
}

// Resolve :
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	return obj, ok
}
