package compiler

import "testing"

// TestDefine :
func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": Symbol{
			Name:  "a",
			Scope: GlobalScope,
			Index: 0,
		},
		"b": Symbol{
			Name:  "b",
			Scope: GlobalScope,
			Index: 1,
		},
		"c": Symbol{
			Name:  "c",
			Scope: LocalScope,
			Index: 0,
		},
		"d": Symbol{
			Name:  "d",
			Scope: LocalScope,
			Index: 1,
		},
		"e": Symbol{
			Name:  "e",
			Scope: LocalScope,
			Index: 0,
		},
		"f": Symbol{
			Name:  "f",
			Scope: LocalScope,
			Index: 1,
		},
	}

	global := InitializeSymbolTable()

	a := global.Define("a")

	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	b := global.Define("b")

	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}

	firstLocal := InitializeEnclosedSymbolTable(global)

	c := firstLocal.Define("c")

	if c != expected["c"] {
		t.Errorf("expected a=%+v, got=%+v", expected["c"], c)
	}

	d := firstLocal.Define("d")

	if d != expected["d"] {
		t.Errorf("expected a=%+v, got=%+v", expected["d"], d)
	}

	secondLocal := InitializeEnclosedSymbolTable(firstLocal)

	e := secondLocal.Define("e")

	if e != expected["e"] {
		t.Errorf("expected a=%+v, got=%+v", expected["e"], e)
	}

	f := secondLocal.Define("f")

	if f != expected["f"] {
		t.Errorf("expected a=%+v, got=%+v", expected["f"], f)
	}
}

// TestResolveGlobal :
func TestResolveGlobal(t *testing.T) {
	global := InitializeSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		Symbol{
			Name:  "a",
			Scope: GlobalScope,
			Index: 0,
		},
		Symbol{
			Name:  "b",
			Scope: GlobalScope,
			Index: 1,
		},
	}

	for _, symbol := range expected {
		result, ok := global.Resolve(symbol.Name)

		if !ok {
			t.Errorf("name %s not resolvable", symbol.Name)

			continue
		}

		if result != symbol {
			t.Errorf("expected '%s' to resolve to %+v, got=%+v", symbol.Name, symbol, result)
		}
	}
}

// TestResolveLocal :
func TestResolveLocal(t *testing.T) {
	global := InitializeSymbolTable()

	global.Define("a")
	global.Define("b")

	local := InitializeEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		Symbol{
			Name:  "a",
			Scope: GlobalScope,
			Index: 0,
		},
		Symbol{
			Name:  "b",
			Scope: GlobalScope,
			Index: 1,
		},
		Symbol{
			Name:  "c",
			Scope: LocalScope,
			Index: 0,
		},
		Symbol{
			Name:  "d",
			Scope: LocalScope,
			Index: 1,
		},
	}

	for _, symbol := range expected {
		result, ok := local.Resolve(symbol.Name)

		if !ok {
			t.Errorf("name %s not resolvable", symbol.Name)

			continue
		}

		if symbol != result {
			t.Errorf("expected %s to resolve to %+v, got=%+v", symbol.Name, symbol, result)
		}
	}
}

// TestResolveNestedLocal :
func TestResolveNestedLocal(t *testing.T) {
	global := InitializeSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := InitializeEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := InitializeEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				Symbol{
					Name:  "a",
					Scope: GlobalScope,
					Index: 0,
				},
				Symbol{
					Name:  "b",
					Scope: GlobalScope,
					Index: 1,
				},
				Symbol{
					Name:  "c",
					Scope: LocalScope,
					Index: 0,
				},
				Symbol{
					Name:  "d",
					Scope: LocalScope,
					Index: 1,
				},
			},
		},
		{
			secondLocal,
			[]Symbol{
				Symbol{
					Name:  "a",
					Scope: GlobalScope,
					Index: 0,
				},
				Symbol{
					Name:  "b",
					Scope: GlobalScope,
					Index: 1,
				},
				Symbol{
					Name:  "e",
					Scope: LocalScope,
					Index: 0,
				},
				Symbol{
					Name:  "f",
					Scope: LocalScope,
					Index: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		for _, symbol := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(symbol.Name)

			if !ok {
				t.Errorf("name '%s' not resolvable", symbol.Name)

				continue
			}

			if symbol != result {
				t.Errorf("expected '%s' to resolve to %+v, got=%+v", symbol.Name, symbol, result)
			}
		}
	}
}

// TestDefineResolveBuiltins :
func TestDefineResolveBuiltins(t *testing.T) {
	global := InitializeSymbolTable()
	firstLocal := InitializeEnclosedSymbolTable(global)
	secondLocal := InitializeEnclosedSymbolTable(firstLocal)
	expected := []Symbol{
		Symbol{
			Name:  "a",
			Scope: BuiltinScope,
			Index: 0,
		},
		Symbol{
			Name:  "c",
			Scope: BuiltinScope,
			Index: 1,
		},
		Symbol{
			Name:  "e",
			Scope: BuiltinScope,
			Index: 2,
		},
		Symbol{
			Name:  "f",
			Scope: BuiltinScope,
			Index: 3,
		},
	}

	for index, value := range expected {
		global.DefineBuiltin(index, value.Name)
	}

	for _, table := range []*SymbolTable{
		global,
		firstLocal,
		secondLocal,
	} {
		for _, symbol := range expected {
			result, ok := table.Resolve(symbol.Name)

			if !ok {
				t.Errorf("name %s not resolvable", symbol.Name)

				continue
			}

			if symbol != result {
				t.Errorf("expected %s to resolve to %+v, got= %+v", symbol.Name, symbol, result)
			}
		}
	}
}
