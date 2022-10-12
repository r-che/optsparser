package optsparser

import (
	"testing"
	"os"
	"time"
)

const (
	stubApp	=	"test-optsparser-app"
)

// Disallow Usage() do os.Exit
func init() {
	usageDoExit = false
}

func TestParser(t *testing.T) {
	// Create new parser
	p := NewParser(stubApp,	// application name
		// Required arguments
		"strval-required",
		"duration-value",
		"intval-required",
	)

	// General description of application
	p.SetGeneralDescr("\n$ test-app --REQUIRED-KEYS... [--optional-keys ...]\n")

	//
	// Add options
	//

	// Add separator - title of options group
	p.AddSeparator(">> Boolean parameters")
	var yesno bool
	p.AddBool("yesno|y", "some boolean value", &yesno, true)

	p.AddSeparator("")
	p.AddSeparator(">> String-based parameters")

	var strVal string
	p.AddString("strval-required|s", "some required string value", &strVal, "")
	var strVal2 string
	p.AddString("strval-def-empty|S", "string value with empty default", &strVal2, "")
	var strVal3 string
	p.AddString("strval", "some string value with defaults", &strVal3, "default string")

	var durationVal time.Duration
	p.AddDuration("duration-value|D", "some duration data", &durationVal, 0)

	p.AddSeparator("")
	p.AddSeparator(">> Integer-based parameters")

	var intVal int
	p.AddInt("intval-required|i", "some integer value", &intVal, -10)

	var int64Val int64
	p.AddInt64("int64val", "some integer64 value", &int64Val, -100)

	var uintVal uint
	p.AddUint("uintval", "some unsigned integer value", &uintVal, 10)

	var uint64Val uint64
	p.AddUint64("uint64val", "some unsigned integer64 value", &uint64Val, 100)

	p.AddSeparator("")
	p.AddSeparator(">> Float64-based parameters")

	var floatVal float64
	p.AddFloat64("floatval", "some float value", &floatVal, 0.0)

	os.Args = os.Args[:1]
	p.Parse()
}
