package optsparser

import (
	"fmt"
	"time"
)

func Example_parse() {
	const exampleAppName = "test-app"

	// XXX Assumed that os.Args contains a list of flags, values and arguments
	// Create new parser
	p := NewParser(exampleAppName,
		// List of required options
		"strval-required",
		"duration-value",
		"intval",
	).
		SetGeneralDescr("\n$ " + exampleAppName + " --REQUIRED-KEYS ... [--optional-keys ...] [ARGUMENTS ...]\n")

	// Define options
	type optsType struct {
		vBool bool
		vStrRq, vStrDef string
		vDur time.Duration
		vInt int; vInt64 int64
		vUInt uint; vUInt64 uint64
		vFloat float64
	}
	opts := optsType{}

	// Add separator as a title of an option group
	p.AddSeparator(">> Boolean parameters")
	p.AddBool("yesno|y", "some boolean value", &opts.vBool, true)

	p.AddSeparator(
		"", // Add empty line to break usage output
		">> String-based parameters",
		">> One required and two parameters with defaults are supported",
	)
	p.AddString("strval-required|S", "some required string value", &opts.vStrRq, "")
	p.AddString("strval-default|s", "some string value with defaults", &opts.vStrDef, "default string")
	p.AddDuration("duration-value|D", "some duration data", &opts.vDur, 0)

	p.AddSeparator("",	// empty line
		">> Integer-based parameters")
	p.AddInt("intval|i", "some integer value", &opts.vInt, -10)
	p.AddInt64("int64val", "some integer64 value", &opts.vInt64, -100)
	p.AddUint("uintval", "some unsigned integer value", &opts.vUInt, 10)
	p.AddUint64("uint64val", "some unsigned integer64 value", &opts.vUInt64, 100)

	p.AddSeparator( "", // empty line
		">> Float64-based parameters")
	p.AddFloat64("floatval", "some float value", &opts.vFloat, 0.0)

	//nolint:errcheck // Perform options parsing, if something is wrong - Parse causes exit the program
	p.Parse()

	// By default this point is reached only if the parsing is successful - no additional checks are needed,
	// but you can hadle the parsing error at your own discretion, see SetUsageOnFail for details
	fmt.Printf("Options are: %#v\n", opts)
	fmt.Printf("Arguments are: %#v\n", p.Args())
}

func Example_handleParseError() {
	// Create new parser
	p := NewParser("test parser",
		"strval-required",
	).
		SetGeneralDescr("\n$ test-app --required-keys ... [--optional-keys ...]\n").
		SetUsageOnFail(false)

	var strVal string
	p.AddString("strval-required|s", "", &strVal, "")

	if err := p.Parse(); err != nil {
		// Write your own error handler here
		fmt.Printf("optsparser.Parse() returned an error: %v\n", err)
		return
	}

	// Options successfully parsed, write the main program code below
}
