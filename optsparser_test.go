package optsparser

import (
	"fmt"
	"bytes"
	"testing"
	"os"
	"io"
	"sort"
)

const (
	stubApp	=	"test-optsparser-app"
)

// Disallow Usage() do os.Exit
func init() {
	usageDoExit = false
}

func TestParser(t *testing.T) {
	// Save current value of os.Args because it will be replaced by test values
	origArgs := make([]string, 0, len(os.Args))
	copy(origArgs, os.Args)
	// Get the binary name
	binName := os.Args[0]
	// Recover on exiting from function
	defer func() {
		os.Args = origArgs
	}()

	// Get tests names and sort them
	names := make([]string, 0, len(tests))
	for name := range tests {
		names = append(names, name)
	}
	sort.Strings(names)

	// Run tests sorted by names
	for _, testN := range names {
		// Get test
		test := tests[testN]

		// Reset usage triggered flag
		usageTriggered = false

		// Make a buffer to catch parser's output
		tOut := &bytes.Buffer{}

		// Create new parser
		p := NewParser(stubApp,	// application name
			test.required...,
		).SetOutput(tOut)

		// Function to automate adding separators
		sepN := -1
		sep := func() {
			if v, ok := test.separators[sepN]; ok {
				p.AddSeparator(v)
			}
			sepN++
		}
		to := testOpts{}

		// Function to automate selection of option names
		opt := func(ot, def string) string {
			if test.keys == nil {
				return def
			}
			if v, ok := test.keys[ot]; ok {
				return v
			}
			return def
		}

		sep()	//	[-1] some kind of additional description
		sep()	//	[0]
		p.AddBool(opt("bool", "bool-opt|b"), "boolean value", &to.vBool, test.defaults.vBool)
		sep()	//	[1]
		p.AddString(opt("string", "string-opt|s"), "string value", &to.vString, test.defaults.vString)
		sep()	//	[2]
		p.AddInt(opt("int", "int-opt|i"), "int value", &to.vInt,  test.defaults.vInt)
		sep()	//	[3]
		p.AddInt64(opt("int64", "int64-opt|I"), "int64 value", &to.vInt64, test.defaults.vInt64)
		sep()	//	[4]
		p.AddFloat64(opt("float64", "float64-opt|f"), "float64 value", &to.vFloat64, test.defaults.vFloat64)
		sep()	//	[5]
		p.AddDuration(opt("duration", "duration-opt|d"), "duration value", &to.vDuration, test.defaults.vDuration)
		sep()	//	[6]
		p.AddUint(opt("uint", "uint-opt|u"), "uint value", &to.vUint, test.defaults.vUint)
		sep()	//	[7]
		p.AddUint64(opt("uint64", "uint64-opt|U"), "uint64 value", &to.vUint64, test.defaults.vUint64)
		sep()	//	[8]
		p.AddVar(opt("var", "var-ymd-opt|V"), "var value", &to.vVar)
		sep()	//	[9]


		// Replace real command arguments
		os.Args = append([]string{binName}, test.args...)

		// Do parsing
		p.Parse()

		// Is test should be OK?
		if test.needOK {
			// Is it true?
			if usageTriggered {
				// False, Usage called
				t.Errorf("%q parse failed: args - %#v, test output:" +
					"\n-------- Start output --------\n%s\n-------- End output --------",
					testN, test.args, tOut.String())
				continue
			}

			// Need to compare parsed and expected results
			if to != test.want {
				t.Errorf("%q incorrect Parse result: want - %#v got - %#v, args - %#v", testN, test.want, to, test.args)
				continue
			}

			// Success, run next test
			continue
		}

		// Test should be failed
		if !usageTriggered {
			t.Errorf("%q incorrect Parse result - test should fail but it succeeds; got - %#v, args - %#v",
				testN, to, test.args)
		}

		// Success, test failed as expected
	}
}

func TestAddIncorrect(t *testing.T) {
	//
	// Test cases
	//
	var bVal bool

	type panicArgs struct {
		opt, usage string
		def bool
	}

	tests := [][]panicArgs {
		// Short and long options should not be the same
		{ { `opt1|opt1`, `boolean value #1 - short and long options are the same`, false } },

		// Too long short option
		{ { `opt2|not-so-short`, `boolean value #2 - invalid length of short option`, false } },
		// Empty short option
		{ { `opt3|`, `boolean value #3 - invalid length of short option`, false } },
		// Empty long option with short
		{ { `|o`, `boolean value #4 - invalid length of short option`, false } },

		// Redefining of existing long-option
		{
			{ `opt5`, `boolean value #4.0 - redefine long opt`, false },
			{ `opt5`, `boolean value #4.1 - redefine long opt`, true },
		},

		// Redefining of existing short-option
		{
			{ `opt6|O`, `boolean value #5.0 - redefine short opt`, false },
			{ `opt7|O`, `boolean value #5.1 - redefine short opt`, true },
		},
	}

	// Function to handle panic
	testPanic := func(p *OptsParser, argsSet []panicArgs) (err error) {
		// Handle panic
		defer func() {
			if p := recover(); p != nil {
				// Ok, panic as expected
				return
			}

			// No panic, this should not be!
			err = fmt.Errorf(`set %#v did not cause a panic, but it must!`, argsSet)
		}()

		// Run AddBool function with all set of arguments
		for _, args := range argsSet {
			p.AddBool(args.opt, args.usage, &bVal, args.def)
		}

		return nil
	}

	// Run tests from set
	for i, test := range tests {
		// Create new parser
		p := NewParser(stubApp).SetOutput(&bytes.Buffer{})

		if err := testPanic(p, test); err != nil {
			t.Errorf(`[%d] panic case return error: %v`, i, err)
		}
	}
}

func TestAddUnsupported(t *testing.T) {
	// Handle panic
	defer func() {
		switch p := recover(); p.(type) {
		case nil:
			// No panic, this should not be!
			t.Errorf(`adding of unsupported option type did not cause a panic, but it must!`)
		case OptsPanic:
			// Ok, panic as expected
		default:
			t.Errorf("raised unexpected panic: %v", p)
		}
	}()

	// Create new parser
	p := NewParser(stubApp).SetOutput(&bytes.Buffer{})

	// Add incorrect option type
	p.addOpt(`TEST-UNSUPPORTED-TYPE`, `unsup-name`, `unsupported value usage`, nil, nil)
}

func TestRequiredNotAdded(t *testing.T) {
	// Handle panic
	defer func() {
		switch p := recover(); p.(type) {
		case nil:
			// No panic, this should not be!
			t.Errorf(`parsing without adding all required options did not cause a panic, but it must!`)
		case OptsPanic:
			// Ok, panic as expected
		default:
			t.Errorf("raised unexpected panic: %v", p)
		}
	}()

	// Create new parser
	p := NewParser(stubApp,
		`bool-option`,
		`int-option`,
	).SetOutput(&bytes.Buffer{})

	// Add incorrect option type
	p.AddBool(`bool-option`, `required bool option`, new(bool), false)

	// Run parsing without int-option
	p.Parse()
}

//
// Functions required for testing
//
func (p *OptsParser) SetOutput(output io.Writer) *OptsParser {
	p.FlagSet.SetOutput(output)

	return p
}
