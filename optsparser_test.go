package optsparser

import (
	"bytes"
	"testing"
	"os"
	"time"
	"io"
)

const (
	stubApp	=	"test-optsparser-app"
)

// Disallow Usage() do os.Exit
func init() {
	usageDoExit = false
}

type testOpts struct {
	vBool		bool
	vString		string
	vInt		int
	vInt64		int64
	vFloat64	float64
	vDuration	time.Duration
	vUint		uint
	vUint64		uint64
	vVar		any
}

func TestParser(t *testing.T) {
	// Save current value of os.Args because it will be replaced by test values
	origArgs := make([]string, 0, len(os.Args))
	copy(origArgs, os.Args)
	// Recover on exiting from function
	defer func() {
		os.Args = origArgs
	}()

	for testN, test := range tests {
		// Reset usage triggered flag
		usageTriggered = false

		// Make a buffer to catch parser's output
		pOut := &bytes.Buffer{}

		// Create new parser
		p := NewParser(stubApp,	// application name
			// TODO Required arguments
		).SetOutput(pOut)

		to := testOpts{}

		p.AddBool("bool-opt", "boolean value", &to.vBool, test.defaults.vBool)
		p.AddString("string-opt", "string value", &to.vString, test.defaults.vString)
		p.AddInt("int-opt", "int value", &to.vInt,  test.defaults.vInt)
		p.AddInt64("int64-opt", "int64 value", &to.vInt64, test.defaults.vInt64)
		p.AddFloat64("float64-opt", "float64 value", &to.vFloat64, test.defaults.vFloat64)
		p.AddDuration("duration-opt", "duration value", &to.vDuration, test.defaults.vDuration)
		p.AddUint("uint-opt", "uint value", &to.vUint, test.defaults.vUint)
		p.AddUint64("uint64-opt", "uint64 value", &to.vUint64, test.defaults.vUint64)
		// TODO p.AddVar("var-opt", "var value", &to.flag.Value)

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
				t.Errorf("[%d] parse failed: usage message - %#v, args - %#v", testN, pOut.String(), test.args)
				continue
			}

			// Need to compare parsed and expected results
			if to != test.want {
				t.Errorf("[%d] incorrect Parse result: want - %#v got - %#v, args - %#v", testN, test.want, to, test.args)
				continue
			}

			// Success, run next test
			continue
		}

		// Test should be failed
		if !usageTriggered {
			t.Errorf("[%d] incorrect Parse result - test should fail but it is successful; got - %#v, args - %#v",
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
