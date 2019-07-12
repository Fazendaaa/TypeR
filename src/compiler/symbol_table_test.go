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

	a := global.Define("a", false)

	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	b := global.Define("b", false)

	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}

	firstLocal := InitializeEnclosedSymbolTable(global)

	c := firstLocal.Define("c", false)

	if c != expected["c"] {
		t.Errorf("expected a=%+v, got=%+v", expected["c"], c)
	}

	d := firstLocal.Define("d", false)

	if d != expected["d"] {
		t.Errorf("expected a=%+v, got=%+v", expected["d"], d)
	}

	secondLocal := InitializeEnclosedSymbolTable(firstLocal)

	e := secondLocal.Define("e", false)

	if e != expected["e"] {
		t.Errorf("expected a=%+v, got=%+v", expected["e"], e)
	}

	f := secondLocal.Define("f", false)

	if f != expected["f"] {
		t.Errorf("expected a=%+v, got=%+v", expected["f"], f)
	}
}

// TestResolveGlobal :
func TestResolveGlobal(t *testing.T) {
	global := InitializeSymbolTable()
	global.Define("a", false)
	global.Define("b", false)

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

	global.Define("a", false)
	global.Define("b", false)

	local := InitializeEnclosedSymbolTable(global)
	local.Define("c", false)
	local.Define("d", false)

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
	global.Define("a", false)
	global.Define("b", false)

	firstLocal := InitializeEnclosedSymbolTable(global)
	firstLocal.Define("c", false)
	firstLocal.Define("d", false)

	secondLocal := InitializeEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e", false)
	secondLocal.Define("f", false)

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

// TestResolveFreeVariables :
func TestResolveFreeVariables(t *testing.T) {
	global := InitializeSymbolTable()
	global.Define("a", false)
	global.Define("b", false)

	firstLocal := InitializeEnclosedSymbolTable(global)
	firstLocal.Define("c", false)
	firstLocal.Define("d", false)

	secondLocal := InitializeEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e", false)
	secondLocal.Define("f", false)

	tests := []struct {
		table                       *SymbolTable
		expectedSymbols             []Symbol
		expectedFreeVariableSymbols []Symbol
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
			[]Symbol{},
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
					Name:  "c",
					Scope: FreeVariableScope,
					Index: 0,
				},
				Symbol{
					Name:  "d",
					Scope: FreeVariableScope,
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
			[]Symbol{
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
	}

	for _, tt := range tests {
		for _, symbol := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(symbol.Name)

			if !ok {
				t.Errorf("name '%s' not resolvable", symbol.Name)

				continue
			}

			if result != symbol {
				t.Errorf("expected '%s' to resolve to %+v, got=%+v", symbol.Name, symbol, result)
			}
		}

		if len(tt.table.FreeVariableSymbol) != len(tt.expectedFreeVariableSymbols) {
			t.Errorf("wrong number of free variables symbols. got=%d, want=%d", len(tt.table.FreeVariableSymbol), len(tt.expectedFreeVariableSymbols))

			continue
		}

		for index, symbol := range tt.expectedFreeVariableSymbols {
			result := tt.table.FreeVariableSymbol[index]

			if result != symbol {
				t.Errorf("wrong free symbol, got=%+v, want=%+v", result, symbol)
			}
		}
	}
}

// TestResolveUnresolvableFreeVariables :
func TestResolveUnresolvableFreeVariables(t *testing.T) {
	global := InitializeSymbolTable()
	global.Define("a", false)

	firstLocal := InitializeEnclosedSymbolTable(global)
	firstLocal.Define("c", false)

	secondLocal := InitializeEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e", false)
	secondLocal.Define("f", false)

	expected := []Symbol{
		Symbol{
			Name:  "a",
			Scope: GlobalScope,
			Index: 0,
		},
		Symbol{
			Name:  "c",
			Scope: FreeVariableScope,
			Index: 0,
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
	}

	for _, symbol := range expected {
		result, ok := secondLocal.Resolve(symbol.Name)

		if !ok {
			t.Errorf("name %s not resolvable", symbol.Name)

			continue
		}

		if result != symbol {
			t.Errorf("expected %s to resolve to %+v, got=%+v", symbol.Name, symbol, result)
		}
	}

	expectedUnresolvable := []string{
		"b",
		"d",
	}

	for _, name := range expectedUnresolvable {
		_, ok := secondLocal.Resolve(name)

		if ok {
			t.Errorf("name %s resolved, but was expected not to", name)
		}
	}
}

// TestDefineAndResolveFunctionName :
func TestDefineAndResolveFunctionName(t *testing.T) {
	global := InitializeSymbolTable()
	global.DefineFunctionName("a")

	expected := Symbol{
		Name:  "a",
		Scope: FunctionScope,
		Index: 0,
	}

	result, ok := global.Resolve(expected.Name)

	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v", expected.Name, expected, result)
	}
}

// TestShadowingFunctionName :
func TestShadowingFunctionName(t *testing.T) {
	global := InitializeSymbolTable()
	global.DefineFunctionName("a")
	global.Define("a", false)

	expected := Symbol{
		Name:  "a",
		Scope: GlobalScope,
		Index: 0,
	}

	result, ok := global.Resolve(expected.Name)

	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v", expected.Name, expected, result)
	}
}
