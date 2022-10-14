package optsparser

import (
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

		sepN := -1
		sep := func() {
			if v, ok := test.separators[sepN]; ok {
				p.AddSeparator(v)
			}
			sepN++
		}
		to := testOpts{}

		sep()	//	[-1] some kind of additional description
		sep()	//	[0]
		p.AddBool("bool-opt|b", "boolean value", &to.vBool, test.defaults.vBool)
		sep()	//	[1]
		p.AddString("string-opt|s", "string value", &to.vString, test.defaults.vString)
		sep()	//	[2]
		p.AddInt("int-opt|i", "int value", &to.vInt,  test.defaults.vInt)
		sep()	//	[3]
		p.AddInt64("int64-opt|I", "int64 value", &to.vInt64, test.defaults.vInt64)
		sep()	//	[4]
		p.AddFloat64("float64-opt|f", "float64 value", &to.vFloat64, test.defaults.vFloat64)
		sep()	//	[5]
		p.AddDuration("duration-opt|d", "duration value", &to.vDuration, test.defaults.vDuration)
		sep()	//	[6]
		p.AddUint("uint-opt|u", "uint value", &to.vUint, test.defaults.vUint)
		sep()	//	[7]
		p.AddUint64("uint64-opt|U", "uint64 value", &to.vUint64, test.defaults.vUint64)
		sep()	//	[8]
		p.AddVar("var-ymd-opt|V", "var value", &to.vVar)
		sep()	//	[9]

		// Update test arguments to insert command name
		test.args = append([]string{os.Args[0]}, test.args...)

		// Replace real command arguments
		os.Args = test.args

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

//
// Functions required for testing
//
func (p *OptsParser) SetOutput(output io.Writer) *OptsParser {
	p.FlagSet.SetOutput(output)

	return p
}
