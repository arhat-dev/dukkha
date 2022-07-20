package tengo

// SymbolScope represents a symbol scope.
type SymbolScope string

// List of symbol scopes
const (
	ScopeGlobal  SymbolScope = "GLOBAL"
	ScopeLocal   SymbolScope = "LOCAL"
	ScopeBuiltin SymbolScope = "BUILTIN"
	ScopeFree    SymbolScope = "FREE"
)

// Symbol represents a symbol in the symbol table.
type Symbol struct {
	Name          string
	Scope         SymbolScope
	Index         int
	LocalAssigned bool // if the local symbol is assigned at least once
}

type SymbolTableStorage[Self any] interface {
	Names() []string
	Get(name string) (*Symbol, bool)
	Add(sym *Symbol)
	// New creates a new instance of the same type
	New() Self
}

var _ SymbolTableStorage[*SimpleSymbolTableStorage] = (*SimpleSymbolTableStorage)(nil)

type SimpleSymbolTableStorage struct {
	Store map[string]*Symbol
}

func NewSimpleSymbolTableStorage() *SimpleSymbolTableStorage {
	return &SimpleSymbolTableStorage{
		Store: make(map[string]*Symbol),
	}
}

// New implements SymbolTableStorage
func (s *SimpleSymbolTableStorage) New() *SimpleSymbolTableStorage {
	return NewSimpleSymbolTableStorage()
}

// Set implements SymbolTableStorage
func (s *SimpleSymbolTableStorage) Add(sym *Symbol) {
	s.Store[sym.Name] = sym
}

// Get implements SymbolTableStorage
func (s *SimpleSymbolTableStorage) Get(name string) (sym *Symbol, ok bool) {
	sym, ok = s.Store[name]
	return
}

// Names implements SymbolTableStorage
func (s *SimpleSymbolTableStorage) Names() (names []string) {
	names = make([]string, len(s.Store))
	i := 0
	for k := range s.Store {
		names[i] = k
		i++
	}

	return
}

// SymbolTable represents a symbol table.
type SymbolTable[S SymbolTableStorage[S]] struct {
	parent         *SymbolTable[S]
	block          bool
	Store          S
	NumDefinition  int
	MaxDefinition  int
	FreeSymbols    []*Symbol
	BuiltinSymbols []Symbol
}

// NewSymbolTable creates a SymbolTable.
func NewSymbolTable[S SymbolTableStorage[S]](store S) *SymbolTable[S] {
	return &SymbolTable[S]{
		Store: store,
	}
}

// Define adds a new symbol in the current scope.
func (t *SymbolTable[S]) Define(symbol *Symbol) {
	symbol.Index = t.nextIndex()
	t.NumDefinition++

	if t.Parent(true) == nil {
		symbol.Scope = ScopeGlobal

		// if symbol is defined in a block of global scope, symbol index must
		// be tracked at the root-level table instead.
		if p := t.parent; p != nil {
			for p.parent != nil {
				p = p.parent
			}
			t.NumDefinition--
			p.NumDefinition++
		}

	} else {
		symbol.Scope = ScopeLocal
	}
	t.Store.Add(symbol)
	t.updateMaxDefs(symbol.Index + 1)
	return
}

// DefineBuiltin adds all symbols for builtin function.
func (t *SymbolTable[S]) BatchDefineAllBuiltin(syms []Symbol) {
	if t.parent != nil {
		t.parent.BatchDefineAllBuiltin(syms)
		return
	}

	t.BuiltinSymbols = syms
	sz := len(syms)
	for i := 0; i < sz; i++ {
		t.Store.Add(&syms[i])
	}
}

// Resolve resolves a symbol with a given name.
func (t *SymbolTable[S]) Resolve(
	name string,
	recur bool,
) (*Symbol, int, bool) {
	symbol, ok := t.Store.Get(name)
	if ok {
		// symbol can be used if
		if symbol.Scope != ScopeLocal || // it's not of local scope, OR,
			symbol.LocalAssigned || // it's assigned at least once, OR,
			recur { // it's defined in higher level
			return symbol, 0, true
		}
	}

	if t.parent == nil {
		return nil, 0, false
	}

	symbol, depth, ok := t.parent.Resolve(name, true)
	if !ok {
		return nil, 0, false
	}
	depth++

	// if symbol is defined in parent table and if it's not global/builtin
	// then it's free variable.
	if !t.block && depth > 0 &&
		symbol.Scope != ScopeGlobal &&
		symbol.Scope != ScopeBuiltin {
		return t.defineFree(symbol), depth, true
	}
	return symbol, depth, true
}

// Fork creates a new symbol table for a new scope.
func (t *SymbolTable[S]) Fork(block bool) *SymbolTable[S] {
	return &SymbolTable[S]{
		Store:  t.Store.New(),
		parent: t,
		block:  block,
	}
}

// Parent returns the outer scope of the current symbol table.
func (t *SymbolTable[S]) Parent(skipBlock bool) *SymbolTable[S] {
	if skipBlock && t.block {
		return t.parent.Parent(skipBlock)
	}
	return t.parent
}

// GetMaxSymbols returns the total number of symbols defined in the scope.
func (t *SymbolTable[S]) GetMaxSymbols() int {
	return t.MaxDefinition
}

// GetFreeSymbols returns free symbols for the scope.
func (t *SymbolTable[S]) GetFreeSymbols() []*Symbol {
	return t.FreeSymbols
}

// GetBuiltinSymbols returns builtin symbols for the scope.
func (t *SymbolTable[S]) GetBuiltinSymbols() []Symbol {
	if t.parent != nil {
		return t.parent.GetBuiltinSymbols()
	}
	return t.BuiltinSymbols
}

// Names returns the name of all the symbols.
func (t *SymbolTable[S]) Names() []string {
	return t.Store.Names()
}

func (t *SymbolTable[S]) nextIndex() int {
	if t.block {
		return t.parent.nextIndex() + t.NumDefinition
	}
	return t.NumDefinition
}

func (t *SymbolTable[S]) updateMaxDefs(numDefs int) {
	if numDefs > t.MaxDefinition {
		t.MaxDefinition = numDefs
	}
	if t.block {
		t.parent.updateMaxDefs(numDefs)
	}
}

func (t *SymbolTable[S]) defineFree(original *Symbol) *Symbol {
	// TODO: should we check duplicates?
	t.FreeSymbols = append(t.FreeSymbols, original)
	symbol := &Symbol{
		Name:  original.Name,
		Index: len(t.FreeSymbols) - 1,
		Scope: ScopeFree,
	}
	t.Store.Add(symbol)
	return symbol
}
