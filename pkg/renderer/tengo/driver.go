package tengo

import (
	"bytes"
	"fmt"
	"io"

	"arhat.dev/pkg/stringhelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/templateutils"

	tengo "github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
)

const DefaultName = "tengo"

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault(name string) dukkha.Renderer { return &Driver{name: name} }

type Driver struct {
	rs.BaseField `yaml:"-"`

	renderer.BaseRenderer `yaml:",inline"`

	name string
}

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{}, _ []dukkha.RendererAttribute,
) (_ []byte, err error) {
	rawData, err = rs.NormalizeRawData(rawData)
	if err != nil {
		return
	}

	var (
		bufScripts [1][]byte
		scripts    [][]byte
	)
	switch t := rawData.(type) {
	case string:
		bufScripts[0] = stringhelper.ToBytes[byte, byte](t)
		scripts = bufScripts[:]
	case []byte:
		bufScripts[0] = t
		scripts = bufScripts[:]
	case []interface{}:
		scripts = make([][]byte, len(t))

		for i, v := range t {
			scripts[i], err = yamlhelper.ToYamlBytes(v)
			if err != nil {
				return nil, fmt.Errorf(
					"renderer.%s: unexpected list item type %T: %w",
					d.name, v, err,
				)
			}
		}
	default:
		return nil, fmt.Errorf(
			"renderer.%s: unsupported input type %T", d.name, rawData,
		)
	}

	var (
		buf     bytes.Buffer
		store   symtabStore
		modules moduleGetterImpl
		symtab  tengo.SymbolTable[*symtabStore]
	)

	symtab.BatchDefineAllBuiltin(tengo.GetAllBuiltinFunctionSymbols())

	for _, src := range scripts {
		symtab.NumDefinition = int(templateutils.FuncID_COUNT)
		symtab.MaxDefinition = symtab.NumDefinition + 1

		err = runScript(rc, &modules, &symtab, src, &buf)
		if err != nil {
			err = fmt.Errorf("renderer.%s: run tengo script: %w", d.name, err)
			return
		}

		store.ClearTmp()
	}

	return buf.Next(buf.Len()), nil
}

var _ tengo.ModuleGetter[*moduleGetterImpl] = (*moduleGetterImpl)(nil)

type moduleGetterImpl struct {
}

// IsNil implements tengo.ModuleGetter
func (g *moduleGetterImpl) IsNil() bool { return g == nil }

// New implements tengo.ModuleGetter
func (g *moduleGetterImpl) New() *moduleGetterImpl {
	return &moduleGetterImpl{}
}

// Get implements tengo.ModuleGetter
func (*moduleGetterImpl) Get(name string) tengo.Importable {
	return nil
}

var _ tengo.SymbolTableStorage[*symtabStore] = (*symtabStore)(nil)

type symtabStore struct {
	templates templateutils.TemplateFuncs

	tmp map[string]*tengo.Symbol
}

// ClearTmp clears all temporary symbols
func (s *symtabStore) ClearTmp() {
	s.tmp = make(map[string]*tengo.Symbol)
}

// Get implements tengo.SymbolTableStorage
func (s *symtabStore) Get(name string) (ret *tengo.Symbol, ok bool) {
	if ret, ok = s.tmp[name]; ok {
		return
	}

	fid := templateutils.FuncNameToFuncID(name)
	if fid != -1 {
		ok = true
		s.templates.GetByID(fid)
	}

	return
}

// Names implements tengo.SymbolTableStorage
func (s *symtabStore) Names() []string {
	return nil
}

// New implements tengo.SymbolTableStorage
func (s *symtabStore) New() *symtabStore {
	return &symtabStore{
		templates: s.templates,
	}
}

// Set implements tengo.SymbolTableStorage
func (s *symtabStore) Add(sym *tengo.Symbol) {

}

// nolint:unparam
func runScript(
	rc dukkha.RenderingContext,
	modules *moduleGetterImpl,
	symtab *tengo.SymbolTable[*symtabStore],
	script []byte,
	stdout *bytes.Buffer,
) (err error) {
	var (
		globals [tengo.GlobalsSize]tengo.Object
		lines   [1024]int
	)

	srcFile := parser.SourceFile{
		Name:  "(main)",
		Base:  0,
		Size:  len(script),
		Lines: lines[:1],
	}

	parsed, err := parser.NewParser(&srcFile, script, nil).ParseFile()
	if err != nil {
		return
	}

	compiler := tengo.NewCompiler(
		&srcFile,
		symtab,
		globals[:0],
		modules,
		io.Discard,
	)
	compiler.EnableFileImport(false)
	compiler.SetImportDir("")

	err = compiler.Compile(parsed)
	if err != nil {
		return
	}

	// remove duplicates from constants
	vm := tengo.NewVM(compiler.Bytecode(), globals[:symtab.GetMaxSymbols()+1], -1)

	return vm.Run()
}
