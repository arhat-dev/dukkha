package transform

import (
	"strings"
	"unicode/utf8"

	"arhat.dev/pkg/stringhelper"
	"arhat.dev/rs"
	"github.com/benhoyt/goawk/interp"
	"github.com/benhoyt/goawk/parser"
)

type awkSpec struct {
	rs.BaseField

	// Script is the awk script
	Script string `yaml:"script"`

	// CSVInput config used when input mode is csv
	CSVInput *CSVOptions `yaml:"csv_input"`

	// CSVOutput config used when output mode is csv
	CSVOutput *CSVOptions `yaml:"csv_output"`

	// TODO: add variables support
}

type CSVOptions struct {
	rs.BaseField

	// Separator mark character
	Sep string `yaml:"sep"`

	// Comment mark character
	//
	// NOTE: this field is unused in output options
	Comment string `yaml:"comment"`

	// FirstLineHeader treats the first line as header to give each column a name
	//
	// NOTE: this field is unused in output options
	FirstLineHeader bool `yaml:"first_line_header"`
}

func (opts *CSVOptions) InputConfig() (ret interp.CSVInputConfig) {
	if len(opts.Sep) != 0 {
		ret.Separator, _ = utf8.DecodeRuneInString(opts.Sep)
	}

	if len(opts.Comment) != 0 {
		ret.Comment, _ = utf8.DecodeRuneInString(opts.Comment)
	}

	ret.Header = opts.FirstLineHeader

	return
}

func (opts *CSVOptions) OutputConfig() (ret interp.CSVOutputConfig) {
	if len(opts.Sep) != 0 {
		ret.Separator, _ = utf8.DecodeRuneInString(opts.Sep)
	}

	return
}

func (s *awkSpec) Run(rc extendedUserFacingRenderContext, value string) (ret string, err error) {
	var (
		prog  *parser.Program
		input strings.Reader
		sb    strings.Builder
	)
	config := parser.ParserConfig{
		DebugWriter: rc.Stderr(),
	}

	prog, err = parser.ParseProgram(stringhelper.ToBytes[byte, byte](s.Script), &config)
	if err != nil {
		return
	}

	// TODO: cache target
	target, err := interp.New(prog)
	if err != nil {
		return
	}

	inputMode := interp.DefaultMode
	var csvInput interp.CSVInputConfig
	if s.CSVInput != nil {
		inputMode = interp.CSVMode
		csvInput = s.CSVInput.InputConfig()
	}

	outputMode := interp.DefaultMode
	var csvOutput interp.CSVOutputConfig
	if s.CSVOutput != nil {
		outputMode = interp.CSVMode
		csvOutput = s.CSVOutput.OutputConfig()
	}

	input.Reset(value)
	runConfig := interp.Config{
		Stdin:     &input,
		Output:    &sb,
		Error:     rc.Stderr(),
		Argv0:     "goawk",
		NoArgVars: true,

		// TODO: support var
		Vars: []string{},

		// TODO: add custom funcs
		Funcs: map[string]any{},

		NoExec:       false,
		NoFileWrites: false,
		NoFileReads:  false,

		// TODO: support environ? but we have lazy values
		Environ: []string{},

		InputMode:  inputMode,
		CSVInput:   csvInput,
		OutputMode: outputMode,
		CSVOutput:  csvOutput,
	}

	_, err = target.Execute(&runConfig)
	ret = sb.String()

	return
}
