package optsparser
import (
	"flag"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
	"sort"
)

const lsJoinDefault = ", "

var (
	// Auxiliary variable to avoid tests termination on Usage() function.
	usageDoExit		=	true	//nolint:gochecknoglobals // need for testing
	// Auxiliary variable to show that Usage() was triggered.
	usageTriggered	=	false	//nolint:gochecknoglobals // need for testing
)

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
	usageOnFail		bool
}

func NewParser(name string, required ...string) *OptsParser {
	parser := &OptsParser{
		FlagSet:		*flag.NewFlagSet(name, flag.ContinueOnError),
		shToLong:		map[string]string{},
		longOpts:		map[string]*optDescr{},
		orderedList:	[]string{},
		required:		map[string]bool{},
		lsJoinStr:		lsJoinDefault,
		usageOnFail:	true,
	}

	// Set stub to FlagSet.Usage to suppress default output.
	parser.FlagSet.Usage = func() {}

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

func (p *OptsParser) AddBool(optName, usage string, val *bool, dfltVal bool) {
	long, short, shOk := p.parseOptName(typeBool, optName, usage)
	p.BoolVar(val, long, dfltVal, usage)
	if shOk {
		p.BoolVar(val, short, dfltVal, usage)
	}
}

func (p *OptsParser) AddString(optName, usage string, val *string, dfltVal string) {
	long, short, shOk := p.parseOptName(typeString, optName, usage)
	p.StringVar(val, long, dfltVal, usage)
	if shOk {
		p.StringVar(val, short, dfltVal, usage)
	}
}

func (p *OptsParser) AddInt(optName, usage string, val *int, dfltVal int) {
	long, short, shOk := p.parseOptName(typeInt, optName, usage)
	p.IntVar(val, long, dfltVal, usage)
	if shOk {
		p.IntVar(val, short, dfltVal, usage)
	}
}

func (p *OptsParser) AddInt64(optName, usage string, val *int64, dfltVal int64) {
	long, short, shOk := p.parseOptName(typeInt64, optName, usage)
	p.Int64Var(val, long, dfltVal, usage)
	if shOk {
		p.Int64Var(val, short, dfltVal, usage)
	}
}

func (p *OptsParser) AddFloat64(optName, usage string, val *float64, dfltVal float64) {
	long, short, shOk := p.parseOptName(typeFloat64, optName, usage)
	p.Float64Var(val, long, dfltVal, usage)
	if shOk {
		p.Float64Var(val, short, dfltVal, usage)
	}
}

func (p *OptsParser) AddDuration(optName, usage string, val *time.Duration, dfltVal time.Duration) {
	long, short, shOk := p.parseOptName(typeDuration, optName, usage)
	p.DurationVar(val, long, dfltVal, usage)
	if shOk {
		p.DurationVar(val, short, dfltVal, usage)
	}
}

func (p *OptsParser) AddUint(optName, usage string, val *uint, dfltVal uint) {
	long, short, shOk := p.parseOptName(typeUint, optName, usage)
	p.UintVar(val, long, dfltVal, usage)
	if shOk {
		p.UintVar(val, short, dfltVal, usage)
	}
}

func (p *OptsParser) AddUint64(optName, usage string, val *uint64, dfltVal uint64) {
	long, short, shOk := p.parseOptName(typeUint64, optName, usage)
	p.Uint64Var(val, long, dfltVal, usage)
	if shOk {
		p.Uint64Var(val, short, dfltVal, usage)
	}
}

func (p *OptsParser) AddVar(optName, usage string, val flag.Value) {
	long, short, shOk := p.parseOptName(typeVal, optName, usage)
	p.Var(val, long, usage)
	if shOk {
		p.Var(val, short, usage)
	}
}

func (p *OptsParser) AddSeparator(separators ...string) {
	for _, separator := range separators {
		// Skip al returned values
		_, _, _ = p.parseOptName(typeSeparator, "", separator)
		// Add new separator
		sep := p.nextSep()
		// Replace last item of ordered list by separator value
		p.orderedList[len(p.orderedList)-1] = sep
		p.StringVar(new(string), sep, "", separator)
	}
}

func (p *OptsParser) Parse() error {
	// Check for all required options was set by Add...() functions
	for opt, required := range p.required {
		if required {
			doPanic("Option %q is required but not added to parser using Add...() method", opt)
		}
	}

	//
	// Suppress parser's output to avoid duplicate error messages
	//

	// Get current output from parser
	out := p.FlagSet.Output()
	// Redirect output to stub buffer
	p.FlagSet.SetOutput(&bytes.Buffer{})

	// Do parsing
	err := p.FlagSet.Parse(os.Args[1:])

	// Recover the output to allow Usage to print if an error occurs
	// XXX Do not use defer for this call because output
	// XXX has to be recovered BEFORE calling Usage()
	p.FlagSet.SetOutput(out)

	if err != nil {
		// Need to call Usage on fail?
		if p.usageOnFail {
			// Call Usage with error description
			p.Usage(err)
		}

		// Just return error
		return err	//nolint:wrapcheck // Obvious parse error - no need to additional error wrapping
	}

	// Check for all required options were set
	rqSet := p.requiredSet()
	if len(rqSet) == len(p.required) {
		// OK, return no errors
		return nil
	}

	//
	// Some of required options were not set
	//

	// Make sorted list of required options
	opts := make([]string, 0, len(p.required))
	for opt := range p.required {
		opts = append(opts, opt)
	}
	// Sort it
	sort.Strings(opts)

	// List of required options that were not set
	notSet := make([]string, 0, len(p.required))
	for _, opt := range opts {
		if !rqSet[opt] {
			// Required option can be long or short, to print correct
			// number of dashes before the option name use dashes() function
			notSet = append(notSet, dashes(opt) + opt)
		}
	}

	// Create an error
	err = fmt.Errorf("required option(s) is missing: %s", strings.Join(notSet, `, `))

	// Need to call Usage on fail?
	if p.usageOnFail {
		// Call Usage with error description
		p.Usage(err)
	}

	// Otherwise - return error
	return err
}

func (p *OptsParser) parseOptName(optType, optName, usage string) (string, string, bool) {
	// Split option name to long and short
	long, short, shOk := strings.Cut(optName, "|")

	switch {
	case optType == typeSeparator:
		// Skip separator

	// Short and long options should not be the same
	case long == short:
		doPanic("Option of type %q with the usage message %q has inappropriate option name", optType, usage)
	// Check for long option
	case long == "" && shOk:
		// Strange situation, passed something like "|o" as optName
		doPanic(`Invalid specification: %q - if you want to use a short option without long one (e.g. "-%s")` +
			` just use %q as the option name parameter`, optName, short, short)
	// Short should has only one character
	case shOk && len(short) != 1:
		doPanic("Invalid option description %q - length of short option must be == 1", optName)
	}

	// Option description
	p.longOpts[long] = &optDescr{optType: optType}
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

	return long, short, shOk
}

func (p *OptsParser) requiredSet() (map[string]bool) {
	// Check for required options were set
	if len(p.required) == 0 {
		// OK - required options were not set, nothing to check
		return nil
	}
	// List of required options that were set
	rqSet := make(map[string]bool, len(p.required))

	p.Visit(func(f *flag.Flag) {
		// Treat option name as long name
		if _, ok := p.required[f.Name]; ok {
			p.required[f.Name] = true
			// Save this option to map of set options
			rqSet[f.Name] = true
		} else
		// Threat option name as short name
		if _, ok := p.required[p.shToLong[f.Name]]; ok {
			p.required[p.shToLong[f.Name]] = true
			// Save this option to map of set options
			rqSet[p.shToLong[f.Name]] = true
		}
	})

	return rqSet
}

func (p *OptsParser) descrLongOpt(optFlag *flag.Flag) string {
	// Output buffer
	out := bytes.NewBuffer([]byte{})
	// Get long option description
	descr := p.longOpts[optFlag.Name]

	// Value description function
	valDescr := func() string {
		if descr.optType == typeBool {
			// Boolean option
			return "[=true|false]"
		}
		// Option with non-boolean argument
		return fmt.Sprintf(" %s", descr.optType)
	}

	// Is short option exists?
	if short := descr.short; short != "" {
		if p.shortFirst {
			// Print short, join string, then long
			fmt.Fprintf(out, optIndent + "-%s%s" + "%s" + "--%s%s\n",
				short, valDescr(), p.lsJoinStr, optFlag.Name, valDescr())
		} else {
			// Print long, join string, then short
			fmt.Fprintf(out, optIndent + "--%s%s" + "%s" + "-%s%s\n",
				optFlag.Name, valDescr(), p.lsJoinStr, short, valDescr())
		}
	} else {
		// Print only long option name, in fact - long options may be short if only short
		// option was added by p.Add... function, for such case use dashes() function
		// to print correct number of dashes before the option
		fmt.Fprintf(out, optIndent + "%s%s%s\n", dashes(optFlag.Name), optFlag.Name, valDescr())
	}

	// Print usage information
	out.WriteString(helpIndent + optFlag.Usage)

	// Print default value if option is not required
	if _, ok := p.required[optFlag.Name]; ok {
		out.WriteString(" (required option)")
	} else {
		defVal := optFlag.DefValue
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

func (p *OptsParser) SetUsageOnFail(v bool) *OptsParser {
	p.usageOnFail = v

	return p
}

func (p *OptsParser) Usage(errDescr ...error) {
	// Check for custom error description
	if len(errDescr) != 0 {
		fmt.Fprintf(p.Output(), "\nUsage ERROR: %v\n", errDescr[0])
	}

	if name := p.FlagSet.Name(); name == "" {
		fmt.Fprintf(p.Output(), "\nUsage:\n")
	} else {
		fmt.Fprintf(p.Output(), "\nUsage of %s:\n", name)
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

	// XXX This condition will not satisfied only in tests
	if usageDoExit {
		os.Exit(1)
	}
	// XXX For tests purposes - set flag that usage called
	usageTriggered = true
}

func (p *OptsParser) nextSep() string {
	defer func() { p.sepIndex++ }()
	return fmt.Sprintf("%s%d", sepPrefix, p.sepIndex)
}

//
// Auxiliary functions and types
//

// OptsPanic type passed to the panic function to be able to distinguish
// panic produced by package from the panic produced by imported packages.
type OptsPanic string

func doPanic(format string, args ...any) {
	panic(OptsPanic(fmt.Sprintf(format, args...)))
}

func dashes(name string) string {
	// Is it a short name of option?
	if len(name) == 1 {
		// Only one dash required
		return "-"
	}

	// Long name - two dashes required
	return "--"
}
