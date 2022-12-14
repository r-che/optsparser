package optsparser

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"testing"
	"time"
)

const (
	stubApp	=	"test-optsparser-app"
)

// newParser wraps creation of parser to prevent call of os.Exit() inside of Usage() function
func newParser(name string, required ...string) *OptsParser {
	p := NewParser(name, required...)

	// Disallow Usage() do os.Exit()
	p.usageNotExits = true

	return p
}

func TestParser(t *testing.T) {
	t.Parallel()

	// Save current value of os.Args because it will be replaced by test values
	origArgs := make([]string, len(os.Args))
	copy(origArgs, os.Args)
	// Get the binary name
	binName := os.Args[0]
	// Recover on exiting from function
	defer func() {
		os.Args = origArgs
	}()

	// Get tests names and sort them
	names := make([]string, 0, len(parserTests))
	for name := range parserTests {
		names = append(names, name)
	}
	sort.Strings(names)

	// Run tests sorted by names
	for _, testN := range names {
		// XXX Do not run the tests in parallel inside the loop, because
		// each test has to update the shared value of os.Args

		// Get test
		test := parserTests[testN]

		// Make a buffer to catch parser's output
		tOut := &bytes.Buffer{}

		//nolint:varnamelen	// Too obvious case in the test
		// Create new parser
		p := newParser(stubApp,	// application name
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
		to := testOpts{} //nolint:varnamelen	// Too obvious case in the test

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
		p.Parse()	//nolint: errcheck

		// Is test should be OK?
		if test.needOK {
			// Is it true?
			if p.usageTriggered {
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
		if !p.usageTriggered {
			t.Errorf("%q incorrect Parse result - test should fail but it succeeds; got - %#v, args - %#v",
				testN, to, test.args)
		}

		// Success, test failed as expected
	}
}

func TestAddIncorrect(t *testing.T) {
	t.Parallel()

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
			err = fmt.Errorf(`set %#v did not cause a panic, but it must`, argsSet)
		}()

		// Run AddBool function with all set of arguments
		for _, args := range argsSet {
			p.AddBool(args.opt, args.usage, &bVal, args.def)
		}

		return nil
	}

	// Run tests from set
	for i, test := range tests {
		// Capture variables to parallel usage
		i, test := i, test
		t.Run(fmt.Sprintf("panic-test#%d", i), func(t *testing.T) {
			t.Parallel()
			// Create new parser
			p := newParser(stubApp).SetOutput(&bytes.Buffer{})

			if err := testPanic(p, test); err != nil {
				t.Errorf(`[%d] panic case returns error: %v`, i, err)
			}
		})
	}
}

func TestRequiredNotAdded(t *testing.T) {
	t.Parallel()

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
	p := newParser(stubApp,
		`bool-option`,
		`int-option`,
	).SetOutput(&bytes.Buffer{})

	// Add incorrect option type
	p.AddBool(`bool-option`, `required bool option`, new(bool), false)

	// Run parsing without int-option
	p.Parse()	//nolint: errcheck
}

func TestUsage(t *testing.T) {
	t.Parallel()

	// Get new parser
	p, tOut := parserWithPredefinedUsage()

	// Call Usage to get output
	p.Usage(fmt.Errorf("test error for testing usage of %s", stubApp))

	// Compare produced output with expected
	if tOut.String() != expUsageOutputWithErr {
		t.Errorf("output produced by Usage is different from expexted, see below:\n" +
			"\n-------- Want --------\n%s\n" +
			"-------- Got --------\n%s\n",
			expUsageOutputWithErr, tOut.String(),
		)
	}
}

func TestHelpOption(t *testing.T) {
	t.Parallel()

	// Save current value of os.Args because it will be replaced by test values
	origArgs := make([]string, len(os.Args))
	copy(origArgs, os.Args)
	// Recover on exiting from function
	defer func() {
		os.Args = origArgs
	}()

	// Run tests
	for _, args := range [][]string {
		[]string{os.Args[0], "--help"},
		[]string{os.Args[0], "--h"},
	} {
		// Get new parser
		p, tOut := parserWithPredefinedUsage()

		// Replace command line arguments
		os.Args = args

		// Parse options
		if err := p.Parse(); !errors.Is(err, flag.ErrHelp) {
			t.Errorf("Parse returned unexpected error value - %v, want - %v", err, flag.ErrHelp)
			t.FailNow()
		}

		// Check that Usage was not triggered
		if !p.usageTriggered {
			t.Errorf("Usage function was not called by arguments: %v", os.Args)
			t.FailNow()
		}

		// Compare produced output with expected
		if tOut.String() != expUsageOutput {
			t.Errorf("output produced by Usage is different from expexted, see below:\n" +
				"\n-------- Want --------\n%s\n" +
				"-------- Got --------\n%s\n",
				expUsageOutput, tOut.String(),
			)
			t.FailNow()
		}
	}
}

func parserWithPredefinedUsage() (*OptsParser, *bytes.Buffer) {
	// Buffer to save Usage output
	tOut := &bytes.Buffer{}

	//nolint:varnamelen	// Too obvious case in the test
	// Create new parser
	p := newParser(stubApp,
		"strval-required",
		"duration-value",
		"intval",
	).
		SetGeneralDescr("\n$ " + stubApp + " --required-keys ... [--optional-keys ...]\n").
		SetShortFirst(true).
		SetUsageOnFail(false).
		SetLongShortJoinStr(` | `).
		SetOutput(tOut)


	// Add options

	// Add separator - title of option group
	p.AddSeparator(">> Boolean parameters")
	var yesno bool
	p.AddBool("yesno|y", "some boolean value", &yesno, true)

	p.AddSeparator(
		"",	// Add empty line to break usage output
		">> String-based parameters",
		">> One required and two parameters with defaults are supported",
	)

	var strVal string
	p.AddString("strval-required|s", "some required string value", &strVal, "")
	var strVal2 string
	p.AddString("strval-def-empty|S", "string value with empty default", &strVal2, "")
	var strVal3 string
	p.AddString("strval", "some string value with defaults", &strVal3, "default string")

	var durationVal time.Duration
	p.AddDuration("duration-value|D", "some duration data", &durationVal, 0)

	p.AddSeparator(
		"",
		">> Integer-based parameters",
	)

	var intVal int
	p.AddInt("intval|i", "some integer value", &intVal, -10)

	var int64Val int64
	p.AddInt64("int64val", "some integer64 value", &int64Val, -100)

	var uintVal uint
	p.AddUint("uintval", "some unsigned integer value", &uintVal, 10)

	var uint64Val uint64
	p.AddUint64("uint64val", "some unsigned integer64 value", &uint64Val, 100)

	p.AddSeparator(
		"",
		">> Float64-based parameters",
	)

	var floatVal float64
	p.AddFloat64("floatval", "some float value", &floatVal, 0.0)

	return p, tOut
}

func TestUsageNoName(t *testing.T) {
	t.Parallel()

	// Buffer to save Usage output
	tOut := &bytes.Buffer{}
	// Create new parser
	p := newParser("",
		"yesno",
	).
		SetGeneralDescr("\n$ " + stubApp + " --required-keys ... [--optional-keys ...]\n").
		SetOutput(tOut)

	p.AddBool("yesno|y", "some boolean value", new(bool), true)

	p.Usage(fmt.Errorf("test error for testing usage of %s", stubApp))

	// Compare produced output with expected
	if tOut.String() != expUsageNoNameOutput {
		t.Errorf("output produced by Usage(no name) is different from expexted, see below:\n" +
			"\n-------- Want --------\n%s\n" +
			"-------- Got --------\n%s\n",
			expUsageNoNameOutput, tOut.String(),
		)
	}
}

//
// Functions required for testing
//
func (p *OptsParser) SetOutput(output io.Writer) *OptsParser {
	p.FlagSet.SetOutput(output)

	return p
}
