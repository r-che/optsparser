package optsparser
import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type OptsParser struct {
	flag.FlagSet
	shToLong	map[string]string
	longOpts	map[string]*optDescr
	orderedList	[]string
	required	map[string]bool
}

func NewParser(name string, required ...string) *OptsParser {
	parser := &OptsParser{
		FlagSet:		*flag.NewFlagSet(name, flag.ContinueOnError),
		shToLong:		map[string]string{},
		longOpts:		map[string]*optDescr{},
		orderedList:	[]string{},
		required:		map[string]bool{},
	}

	// Set required options
	for _, opt := range required {
		parser.required[opt] = false
	}

	return parser
}

func (p *OptsParser) addOpt(optType, optName, usage string, val, dfltValue interface{}) {
	// Split option name to long and short
	long, short, shOk := strings.Cut(optName, "|")
	switch {
	// Short and long options should not be the same
	case long == short:
		panic("Option of type " + optType + " with usage \"" + usage + "\" has inappropriate option name")
	// Check for long option
	case long == "" && shOk:
		// Replace long by short
		long = short
		// Clear Ok flag to skip short option processing
		shOk = false
	// Short should has only one character
	case shOk && len(short) != 1:
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

	// If this options is required - need to mark it as added to parser
	if _, ok := p.required[long]; ok {
		p.required[long] = true
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
	// Check for all required options was set by Add...() functions
	for opt, val := range p.required {
		if !val {
			panic("Option \"" + opt + "\" is required but not added to parser using Add...() method")
		}
	}

	// Do parsing
	if err := p.FlagSet.Parse(os.Args[1:]); err != nil {
		p.Usage()
	}

	// Check for required options were set
	if len(p.required) != 0 {
		p.Visit(func(f *flag.Flag) {
			// Remove all options that were set from required list

			// Treat option name as long name
			delete(p.required, f.Name)

			// Threat option name as short name
			delete(p.required, p.shToLong[f.Name])
		})

		// Check for some of required options were not set
		if len(p.required) != 0 {
			fmt.Fprintf(os.Stderr, "Some required options are not set:\n")
			for opt := range p.required {
				fmt.Fprintf(os.Stderr, "  --%s\n", opt)
			}
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", p.Name())
			p.Usage()
		}
	}
}

func (p *OptsParser) descrLongOpt(f *flag.Flag) string {
	// Output buffer
	out := strings.Builder{}
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

	// Print default value if option is not required
	if _, ok := p.required[f.Name]; !ok {
		out.WriteString(fmt.Sprintf(" (default: %v)", f.Value))
	} else {
		out.WriteString(" (required option)")
	}

	out.WriteString("\n")

	// Return description
	return out.String()
}

func (p *OptsParser) Usage() {
	for _, opt := range p.orderedList {
		f := p.Lookup(opt)
		// Print option help info
		fmt.Fprint(os.Stderr, p.descrLongOpt(f))
	}
	os.Exit(1)
}
