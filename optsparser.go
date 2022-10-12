package optsparser
import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"bytes"
)

const lsJoinDefault = ", "

// Auxiliary variable to avoid tests termination on Usage() function
var usageDoExit = true

type OptsParser struct {
	flag.FlagSet
	shToLong		map[string]string
	longOpts		map[string]*optDescr
	orderedList		[]string
	required		map[string]bool
	generalDescr	string
	sepIndex		int
	lsJoinStr		string	// long + short join string
	shortFirst		bool
}

func NewParser(name string, required ...string) *OptsParser {
	parser := &OptsParser{
		FlagSet:		*flag.NewFlagSet(name, flag.ContinueOnError),
		shToLong:		map[string]string{},
		longOpts:		map[string]*optDescr{},
		orderedList:	[]string{},
		required:		map[string]bool{},
		lsJoinStr:		lsJoinDefault,
	}

	// Set required options
	for _, opt := range required {
		parser.required[opt] = true
	}

	return parser
}

func (p *OptsParser) SetGeneralDescr(descr string) *OptsParser {
	p.generalDescr = descr

	return p
}

func (p *OptsParser) addOpt(optType, optName, usage string, val, dfltValue interface{}) {
	// Split option name to long and short
	long, short, shOk := strings.Cut(optName, "|")
	switch {
	case optType == typeSeparator:
		// Skip separator

	// Short and long options should not be the same
	case long == short:
		panic("Option of type \"" + optType + "\" with usage \"" + usage + "\" has inappropriate option name")
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
		p.required[long] = false
	}

	// Stub for separators
	var sepStub string

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
	case typeSeparator:
		// Add new separator
		sep := p.nextSep()
		// Replace last item of ordered list by separator value
		p.orderedList[len(p.orderedList)-1] = sep
		p.StringVar(&sepStub, sep, "", usage)
	default:
		panic("Cannot add argument \"" + long + "\" with unsupported type \"" + optType + "\"")
	}
}

func (p *OptsParser) AddBool(optName, usage string, val *bool, dfltVal bool) {
	p.addOpt(typeBool, optName, usage, val, dfltVal)
}

func (p *OptsParser) AddString(optName, usage string, val *string, dfltVal string) {
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

func (p *OptsParser) AddSeparator(text string) {
	p.addOpt(typeSeparator, "", text, nil, nil)
}

func (p *OptsParser) Parse() {
	// Check for all required options was set by Add...() functions
	for opt, required := range p.required {
		if required {
			panic("Option \"" + opt + "\" is required but not added to parser using Add...() method")
		}
	}

	// Do parsing
	if err := p.FlagSet.Parse(os.Args[1:]); err != nil {
		p.Usage()
	}

	// Check for required options were set
	if len(p.required) != 0 {
		// Counter of required option that were set
		nSet := 0
		p.Visit(func(f *flag.Flag) {
			// Remove all options that were set from required list

			// Treat option name as long name
			if _, ok := p.required[f.Name]; ok {
				p.required[f.Name] = true
				nSet++
			} else
			// Threat option name as short name
			if _, ok := p.required[p.shToLong[f.Name]]; ok {
				p.required[p.shToLong[f.Name]] = true
				nSet++
			}

		})

		// Check for some of required options were not set
		if nSet != len(p.required) {
			fmt.Fprintf(p.Output(), "Error, some required options are not set:\n")
			for opt := range p.required {
				fmt.Fprintf(p.Output(), optIndent + "--%s\n", opt)
			}
			fmt.Fprintf(p.Output(), "\nUsage of %s:\n", p.Name())
			p.Usage()
		}
	}
}

func (p *OptsParser) descrLongOpt(f *flag.Flag) string {
	// Output buffer
	out := bytes.NewBuffer([]byte{})
	// Get long option description
	descr := p.longOpts[f.Name]

	// Value description function
	valDescr := func(d *optDescr) string {
		if descr.optType == typeBool {
			// Boolean option
			return "[=true|false]"
		}
		// Option with non-boolean agrument
		return fmt.Sprintf(" %s", descr.optType)
	}

	// Is short option exists?
	if short := descr.short; short != "" {
		if p.shortFirst {
			// Print short, join string, then long
			fmt.Fprintf(out, optIndent + "-%s%s" + "%s" + "--%s%s\n",
				short, valDescr(descr), p.lsJoinStr, f.Name, valDescr(descr))
		} else {
			// Print long, join string, then short
			fmt.Fprintf(out, optIndent + "--%s%s" + "%s" + "-%s%s\n",
				f.Name, valDescr(descr), p.lsJoinStr, short, valDescr(descr))
		}
	} else {
		// Print only long option name
		fmt.Fprintf(out, optIndent + "--%s%s\n", f.Name, valDescr(descr))
	}

	// Print usage information
	out.WriteString(helpIndent + f.Usage)

	// Print default value if option is not required
	if _, ok := p.required[f.Name]; ok {
		out.WriteString(" (required option)")
	} else {
		defVal := f.DefValue
		if defVal == "" {
			// Replace by quotes
			defVal = `""`
		}

		fmt.Fprintf(out, " (default: %v)", defVal)
	}

	out.WriteString("\n")

	// Return description
	return out.String()
}

func (p *OptsParser) SetLongShortJoinStr(join string) *OptsParser {
	p.lsJoinStr = join

	return p
}

func (p *OptsParser) SetShortFirst(v bool) *OptsParser {
	p.shortFirst = v

	return p
}

func (p *OptsParser) Usage(errDescr ...string) {
	// Check for custom error description
	if len(errDescr) != 0 {
		fmt.Fprintf(p.Output(), "\nUsage error: %s\n", errDescr[0])
	}

	// Print common description if set
	if p.generalDescr != "" {
		fmt.Fprintf(p.Output(), "%s\n", p.generalDescr)
	}

	// Reset separators index
	p.sepIndex = 0
	nextSep := p.nextSep()

	for _, opt := range p.orderedList {
		f := p.Lookup(opt)
		// Check is option a separator
		if opt == nextSep {
			// Update the value of the next expected separator
			nextSep = p.nextSep()
			// Print separator, then continue to the next option
			fmt.Fprintf(p.Output(), optIndent + "%s\n", f.Usage)
			continue
		}

		// Print option help info
		fmt.Fprint(p.Output(), p.descrLongOpt(f))
	}

	// XXX This condition is not satisfied only in tests
	if usageDoExit {
		os.Exit(1)
	}
}

func (p *OptsParser) nextSep() string {
	defer func() { p.sepIndex++ }()
	return fmt.Sprintf("%s%d", sepPrefix, p.sepIndex)
}
