package optsparser
import (
	"flag"
	"fmt"
	"os"
	"strings"
	"bytes"
	"time"
)

type OptsParser struct {
	flag.FlagSet
	shToLong	map[string]string
	longOpts	map[string]*optDescr
	orderedList	[]string
	exeName		string
}

func NewParser(name string) *OptsParser {
	parser := &OptsParser{
		FlagSet:		*flag.NewFlagSet(name, flag.ContinueOnError),
		shToLong:		map[string]string{},
		longOpts:		map[string]*optDescr{},
		orderedList:	[]string{},
	}

	// Ignore all outputs produced by standard error functions
	parser.SetOutput(&dummyWriter{})

	return parser
}

func (p *OptsParser) addOpt(optType, optName, usage string, val, dfltValue interface{}) {
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
	case typeBool:
		p.BoolVar(val.(*bool), long, dfltValue.(bool), usage)
		if shOk {
			p.BoolVar(val.(*bool), short, dfltValue.(bool), usage)
		}
	case typeString:
		p.StringVar(val.(*string), long, dfltValue.(string), usage)
		if shOk {
			p.StringVar(val.(*string), short, dfltValue.(string), usage)
		}
	case typeUint:
		p.UintVar(val.(*uint), long, dfltValue.(uint), usage)
		if shOk {
			p.UintVar(val.(*uint), short, dfltValue.(uint), usage)
		}
	case typeUint64:
		p.Uint64Var(val.(*uint64), long, dfltValue.(uint64), usage)
		if shOk {
			p.Uint64Var(val.(*uint64), short, dfltValue.(uint64), usage)
		}
	case typeInt:
		p.IntVar(val.(*int), long, dfltValue.(int), usage)
		if shOk {
			p.IntVar(val.(*int), short, dfltValue.(int), usage)
		}
	case typeInt64:
		p.Int64Var(val.(*int64), long, dfltValue.(int64), usage)
		if shOk {
			p.Int64Var(val.(*int64), short, dfltValue.(int64), usage)
		}
	case typeFloat64:
		p.Float64Var(val.(*float64), long, dfltValue.(float64), usage)
		if shOk {
			p.Float64Var(val.(*float64), short, dfltValue.(float64), usage)
		}
	case typeDuration:
		p.DurationVar(val.(*time.Duration), long, dfltValue.(time.Duration), usage)
		if shOk {
			p.DurationVar(val.(*time.Duration), short, dfltValue.(time.Duration), usage)
		}
	case typeVal:
		p.Var(val.(flag.Value), long, usage)
		if shOk {
			p.Var(val.(flag.Value), short, usage)
		}
	default:
		panic("Cannot add argument \"" + long + "\" with unsupported type \"" + optType + "\"")
	}
}

func (p *OptsParser) AddBool(optName, usage string, val *bool, dfltVal bool) {
	p.addOpt(typeBool, optName, usage, val, dfltVal)
}

func (p *OptsParser) AddStr(optName, usage string, val *string, dfltVal string) {
	p.addOpt(typeString, optName, usage, val, dfltVal)
}

func (p *OptsParser) AddInt(optName, usage string, val *int, dfltVal int) {
	p.addOpt(typeInt, optName, usage, val, dfltVal)
}

func (p *OptsParser) AddInt64(optName, usage string, val *int64, dfltVal int64) {
	p.addOpt(typeInt64, optName, usage, val, dfltVal)
}

func (p *OptsParser) AddFloat64(optName, usage string, val *float64, dfltVal float64) {
	p.addOpt(typeFloat64, optName, usage, val, dfltVal)
}

func (p *OptsParser) AddDuration(optName, usage string, val *time.Duration, dfltVal time.Duration) {
	p.addOpt(typeDuration, optName, usage, val, dfltVal)
}

func (p *OptsParser) AddUint(optName, usage string, val *uint, dfltVal uint) {
	p.addOpt(typeUint, optName, usage, val, dfltVal)
}

func (p *OptsParser) AddUint64(optName, usage string, val *uint64, dfltVal uint64) {
	p.addOpt(typeUint64, optName, usage, val, dfltVal)
}

func (p *OptsParser) AddVar(optName, usage string, val flag.Value) {
	p.addOpt(typeVal, optName, usage, val, nil)
}

func (p *OptsParser) Parse() {
	// Extract programm name
	p.exeName = os.Args[0]

	// Do parsing
	if err := p.FlagSet.Parse(os.Args[1:]); err != nil {
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
		if descr.optType == typeBool {
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

	for _, opt := range p.orderedList {
		f := p.Lookup(opt)
		// Print option help info
		fmt.Print(p.descrLongOpt(f))
	}
}
