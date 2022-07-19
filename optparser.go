package main
import (
	"flag"
	"fmt"
	"os"
	"strings"
	"bytes"
)

func main() {
	fmt.Printf("Provided arguments: %#v\n", os.Args)

	// Create new parser
	p := NewParser("test parser")

	// Add options
	var yesno bool
	p.AddOpt("yesno|y", Bool, "some boolean value", &yesno, true)

	var strVal string
	p.AddOpt("strval|s", String, "some string value", &strVal, "default string")

	var intVal int
	p.AddOpt("intval", Int, "some integer value", &intVal, 10)

	var floatVal float64
	p.AddOpt("floatval", Float, "some float value", &floatVal, 0.0)

	p.Parse()
}

const (
	Bool	=	"bool"
	String	=	"string"
	Int		=	"int"
	Float	=	"float"
)

type optDescr struct {
	// TODO Add required flag
	optType		string
	short		string
}

type OptsParser struct {
	name		string
	shToLong	map[string]string
	longOpts	map[string]*optDescr
	orderedList	[]string
	fs			*flag.FlagSet
	exeName		string
}
type dummyWriter struct {}
func (dw *dummyWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func NewParser(name string) *OptsParser {
	parser := &OptsParser {
		name:			name,
		shToLong:		map[string]string{},
		longOpts:		map[string]*optDescr{},
		orderedList:	[]string{},
	}

	// Create and configure FlagSet object
	parser.fs = flag.NewFlagSet(name, flag.ContinueOnError)
	// Ignore all outputs processed by standard error functions
	parser.fs.SetOutput(&dummyWriter{})

	return parser
}

func (p *OptsParser) AddOpt(optName, optType, usage string, val, dfltValue interface{}) {
	// Split option name to long and short
	long, short, shOk := strings.Cut(optName, "|")
	// Short should has only one character
	if shOk && len(short) != 1 {
		panic("Invalid option description \"" + optName + "\" - length of short option must be == 1")
	}
	// Option description
	descr := optDescr{optType: optType}
	p.longOpts[long] = &descr
	p.orderedList = append(p.orderedList, long)

	// If short option was provided
	if shOk {
		// Set match between short and long options
		p.shToLong[short] = long
		p.longOpts[long].short = short
	}

	// Using standard flag functions
	switch optType {
		case Bool:
			p.fs.BoolVar(val.(*bool), long, dfltValue.(bool), usage)
			if shOk {
				p.fs.BoolVar(val.(*bool), short, dfltValue.(bool), usage)
			}
		case String:
			p.fs.StringVar(val.(*string), long, dfltValue.(string), usage)
			if shOk {
				p.fs.StringVar(val.(*string), short, dfltValue.(string), usage)
			}
		case Int:
			p.fs.IntVar(val.(*int), long, dfltValue.(int), usage)
			if shOk {
				p.fs.IntVar(val.(*int), short, dfltValue.(int), usage)
			}
		case Float:
			p.fs.Float64Var(val.(*float64), long, dfltValue.(float64), usage)
			if shOk {
				p.fs.Float64Var(val.(*float64), short, dfltValue.(float64), usage)
			}
		default:
			panic("Cannot add argument \"" + long + "\" with unsupported type \"" + optType + "\"")
	}
}

func (p *OptsParser) Parse() {
	// Extract programm name
	p.exeName = os.Args[0]

	// Do parsing
	if err := p.fs.Parse(os.Args[1:]); err != nil {
		p.Usage()
	}
}

func (p *OptsParser) descrLongOpt(f *flag.Flag) string {
	// Output buffer
	out := bytes.Buffer{}
	// Get long option description
	descr := p.longOpts[f.Name]

	// Value description function
	valDescr := func(d *optDescr) string {
		if descr.optType == Bool {
			// Boolean option
			return "[=true|false]\n"
		}
		// Option with non-boolean agrument
		return fmt.Sprintf(" %s\n", descr.optType)
	}

	// Is short option exists?
	if short := descr.short; short != "" {
		out.WriteString(fmt.Sprintf("  -%s%s", short, valDescr(descr)))
	}
	// Long option name
	out.WriteString(fmt.Sprintf("  --%s%s", f.Name, valDescr(descr)))

	// Print usage information
	out.WriteString("      " + f.Usage)

	// Print default
	// TODO Process "required" argument to skip default output for required arguments
	out.WriteString(fmt.Sprintf(" (default: %v)", f.Value))

	out.WriteString("\n")

//&flag.Flag{Name:"yesno", Usage:"some boolean value", Value:(*flag.boolValue)(0xc00001a0d8), DefValue:"false"}

	// Description

	return out.String()
}

func (p *OptsParser) Usage() {
	if(p.name != "") {
		fmt.Println("Usage of", p.name + ":" )
	} else {
		fmt.Println("Usage:")
	}

	for _, opt := range p.orderedList {
		f := p.fs.Lookup(opt)
		// Print option help info
		fmt.Print(p.descrLongOpt(f))
	}
}
