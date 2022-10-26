package optsparser

import (
	"fmt"
	"strings"
	"time"
	"strconv"
)

type testOpts struct {
	vBool		bool
	vString		string
	vInt		int
	vInt64		int64
	vFloat64	float64
	vDuration	time.Duration
	vUint		uint
	vUint64		uint64
	vVar		testOptVarYMD
}

// Reference value.
var ref = testOpts{
	vBool:		true,
	vString:	`I think, therefore I am`,
	vInt:		-430,
	vInt64:		-59604644783353249,
	vFloat64:	3.141592,
	vDuration:	time.Hour * 4176,
	vUint:		220414,
	vUint64:	40208000000000,
	vVar:		testOptVarYMD(20221014),
}

// Set of tests.
var tests = map[string]struct{
	keys		map[string]string
	defaults	testOpts
	required	[]string
	args		[]string
	want		testOpts
	separators	map[int]string
	needOK		bool
}{
	//
	// Test OK cases
	//

	// No required options, no special default values
	`00[ok]no-required,no-defaults`: {
		args: []string{
			`--bool-opt`,
			`--string-opt`,		`I think, therefore I am`,
			`--int-opt`,		`-430`,								// Dead Sea level below sea level
			`--int64-opt`,		`-59604644783353249`,				// Leyland prime number
			`--float64-opt`,	`3.141592`,
			`--duration-opt`,	`250560m`,							// 4176 hours => 174 days
			`--uint-opt`,		`220414`,							// The End and the Beginning
			`--uint64-opt`,		`40208000000000`,					// distance between Sun and Proxima Centauri, km
			`--var-ymd-opt`,	`2022.10.14`,
			`command line arg#1`, `command line arg#2`, `command line arg#3`,
		},
		separators: map[int]string {
			-1:	`## Start options list`,
			0:	`# Boolean options`,
			1:	`# String options`,
			2:	`# Integer options`,
			4:	`# Float options`,
			5:	`# Duration options`,
			6:	`# Unsigned int options`,
		},
		want: ref,	// expected result equal reference value
		needOK:	true,
	},

	// No required options, with special default values
	`01[ok]no-required_with-defaults`: {
		defaults: ref,	// pass reference value as defaults
		args: []string{
			`--var-ymd-opt`,	`2022.10.14`,	// XXX AddVar() does not support default value, need to provide through arguments
			`command line arg#1`, `command line arg#2`, `command line arg#3`,
		},
		want: ref,		// expect the same value after parsing
		needOK:	true,
	},

	// Short options test - no required options, no special default values
	`02[ok]short-opts_no-required,no-defaults`: {
		args: []string{
			`-b`,
			`-s`, `I think, therefore I am`,
			`-i`, `-430`,							// Dead Sea level below sea level
			`-I`, `-59604644783353249`,				// Leyland prime number
			`-f`, `3.141592`,
			`-d`, `250560m`,						// 4176 hours => 174 days
			`-u`, `220414`,							// The End and the Beginning
			`-U`, `40208000000000`,					// distance between Sun and Proxima Centauri, km
			`-V`,`2022.10.14`,
		},
		want: ref,	// expected result equal reference value
		needOK:	true,
	},

	// Test required options, no special default values
	`03[ok]required-opts,no-defaults`: {
		required: []string{
			`bool-opt`,
			`int64-opt`,
			`duration-opt`,
			`var-ymd-opt`,
		},
		args: []string{
			`--bool-opt`,
			`--string-opt`,		`I think, therefore I am`,
			`--int-opt`,		`-430`,								// Dead Sea level below sea level
			`--int64-opt`,		`-59604644783353249`,				// Leyland prime number
			`--float64-opt`,	`3.141592`,
			`--duration-opt`,	`250560m`,							// 4176 hours => 174 days
			`--uint-opt`,		`220414`,							// The End and the Beginning
			`--uint64-opt`,		`40208000000000`,					// distance between Sun and Proxima Centauri, km
			`--var-ymd-opt`,	`2022.10.14`,
			// Test work without command line arguments
		},
		want: ref,	// expected result equal reference value
		needOK:	true,
	},

	// Test parsing with only short options
	`04[ok]only-short_with-required,no-defaults`: {
		keys: map[string]string {
			`bool`:		`b`,
			`string`:	`s`,
			`int`:		`i`,
			`int64`:	`I`,
			`float64`:	`f`,
			`duration`:	`d`,
			`uint`:		`u`,
			`uint64`:	`U`,
			`var`:		`V`,
		},
		required: []string{
			`b`,
			`I`,
			`d`,
			`V`,
		},
		args: []string{
			`-b`,
			`-s`,	`I think, therefore I am`,
			`-i`,	`-430`,								// Dead Sea level below sea level
			`-I`,	`-59604644783353249`,				// Leyland prime number
			`-f`,	`3.141592`,
			`-d`,	`250560m`,							// 4176 hours => 174 days
			`-u`,	`220414`,							// The End and the Beginning
			`-U`,	`40208000000000`,					// distance between Sun and Proxima Centauri, km
			`-V`,	`2022.10.14`,
			`command line arg#1`, `command line arg#2`, `command line arg#3`,
		},
		want: ref,	// expected result equal reference value
		needOK:	true,
	},

	// Test parsing with only long options
	`05[ok]only-long_with-required,no-defaults`: {
		keys: map[string]string {
			`bool`:		`bool-opt`,
			`string`:	`string-opt`,
			`int`:		`int-opt`,
			`int64`:	`int64-opt`,
			`float64`:	`float64-opt`,
			`duration`:	`duration-opt`,
			`uint`:		`uint-opt`,
			`uint64`:	`uint64-opt`,
			`var`:		`var-ymd-opt`,
		},
		required: []string{
			`bool-opt`,
			`int64-opt`,
			`duration-opt`,
			`var-ymd-opt`,
		},
		args: []string{
			`--bool-opt`,
			`--string-opt`,		`I think, therefore I am`,
			`--int-opt`,		`-430`,								// Dead Sea level below sea level
			`--int64-opt`,		`-59604644783353249`,				// Leyland prime number
			`--float64-opt`,	`3.141592`,
			`--duration-opt`,	`250560m`,							// 4176 hours => 174 days
			`--uint-opt`,		`220414`,							// The End and the Beginning
			`--uint64-opt`,		`40208000000000`,					// distance between Sun and Proxima Centauri, km
			`--var-ymd-opt`,	`2022.10.14`,
			`command line arg#1`, `command line arg#2`, `command line arg#3`,
		},
		want: ref,	// expected result equal reference value
		needOK:	true,
	},

	// No required options, no special default values
	`06[ok]required-long_provided-short`: {
		required: []string{
			`bool-opt`,
			`int64-opt`,
			`duration-opt`,
			`var-ymd-opt`,
		},
		args: []string{
			`-b`,
			`-s`,	`I think, therefore I am`,
			`-i`,	`-430`,								// Dead Sea level below sea level
			`-I`,	`-59604644783353249`,				// Leyland prime number
			`-f`,	`3.141592`,
			`-d`,	`250560m`,							// 4176 hours => 174 days
			`-u`,	`220414`,							// The End and the Beginning
			`-U`,	`40208000000000`,					// distance between Sun and Proxima Centauri, km
			`-V`,	`2022.10.14`,
			`command line arg#1`, `command line arg#2`, `command line arg#3`,
		},
		want: ref,	// expected result equal reference value
		needOK:	true,
	},

	//
	// Test FAIL cases
	//

	// Incorrect option passed
	`10[fail]incorrect-option`: {
		args: []string{
			`--bool-opt`,
			`--string-opt`,		`I think, therefore I am`,
			`--int-opt`,		`-430`,								// Dead Sea level below sea level
			`--int64-opt`,		`-59604644783353249`,				// Leyland prime number
			`--float64-opt`,	`3.141592`,
			`--special-incorrect-option`,
			`--duration-opt`,	`250560m`,							// 4176 hours => 174 days
			`--uint-opt`,		`220414`,							// The End and the Beginning
			`--uint64-opt`,		`40208000000000`,					// distance between Sun and Proxima Centauri, km
			`--var-ymd-opt`,	`2022.10.14`,
			`command line arg#1`, `command line arg#2`, `command line arg#3`,
		},
		want: ref,	// expected result equal reference value
		needOK:	false,
	},

	// Incorrect values passed - all Parse have to fail
	`11.0[fail]incorrect-option-bool`:		{ args: []string{`--bool-opt=invalid`}, want: testOpts{} },
	// XXX Skip testing the string type because it is impossible to pass something inappropriate for the string type
	`11.1[fail]incorrect-option-int`:			{ args: []string{`--int-opt`, `f430`}, want: testOpts{} },
	`11.2[fail]incorrect-option-int64`:		{ args: []string{`--int64-opt`, `59604644783353249.1`}, want: testOpts{} },
	`11.3[fail]incorrect-option-float64`:		{ args: []string{`--float64-opt`, `3.G+e0`, }, want: testOpts{} },
	`11.4[fail]incorrect-option-duration`:	{ args: []string{`--duration-opt`, `1Y`}, want: testOpts{} },
	`11.5[fail]incorrect-option-uint`:		{ args: []string{`--uint-opt`, `-220414`}, want: testOpts{} },
	`11.6[fail]incorrect-option-uint64`:		{ args: []string{`--uint64-opt`, `-40208000000000`}, want: testOpts{} },
	`11.7[fail]incorrect-option-var`:			{ args: []string{`--var-ymd-opt`, `2022/10/14`, }, want: testOpts{} },

	// Test required option missing
	`12[fail]required_short-opts`: {
		required: []string{
			`bool-opt`,
			`int64-opt`,
			`duration-opt`,
			`var-ymd-opt`,
		},
		args: []string{
			`--bool-opt`,
			`--string-opt`,		`I think, therefore I am`,
			`--int-opt`,		`-430`,								// Dead Sea level below sea level
			`--int64-opt`,		`-59604644783353249`,				// Leyland prime number
			`--float64-opt`,	`3.141592`,
			`--uint-opt`,		`220414`,							// The End and the Beginning
			`--uint64-opt`,		`40208000000000`,					// distance between Sun and Proxima Centauri, km
			`--var-ymd-opt`,	`2022.10.14`,
			`command line arg#1`, `command line arg#2`, `command line arg#3`,
		},
		want: ref,	// expected result equal reference value
		needOK:	false,
	},

	// Test short required options missing
	`13[fail]required_short-opts`: {
		keys: map[string]string {
			`bool`:		`b`,
			`int64`:	`I`,
			`duration`:	`d`,
			`var`:		`V`,
		},
		required: []string{
			`b`,
			`I`,
			`d`,
			`V`,
		},
		args: []string{
			`-s`,	`I think, therefore I am`,
			`-i`,	`-430`,								// Dead Sea level below sea level
			`-I`,	`-59604644783353249`,				// Leyland prime number
			`-f`,	`3.141592`,
			`-u`,	`220414`,							// The End and the Beginning
			`-U`,	`40208000000000`,					// distance between Sun and Proxima Centauri, km
			`command line arg#1`, `command line arg#2`, `command line arg#3`,
		},
		want: ref,	// expected result equal reference value
		needOK:	false,
	},
}

// Expected Usage outputs.
const expUsageOutput = `
Usage ERROR: test error for testing usage of ` +  stubApp + `

Usage of ` + stubApp + `:

$ ` + stubApp + ` --required-keys ... [--optional-keys ...]

    >> Boolean parameters
    -y[=true|false] | --yesno[=true|false]
      some boolean value (default: true)
    
    >> String-based parameters
    >> One required and two parameters with defaults are supported
    -s string | --strval-required string
      some required string value (required option)
    -S string | --strval-def-empty string
      string value with empty default (default: "")
    --strval string
      some string value with defaults (default: default string)
    -D duration | --duration-value duration
      some duration data (required option)
    
    >> Integer-based parameters
    -i int | --intval int
      some integer value (required option)
    --int64val int64
      some integer64 value (default: -100)
    --uintval uint
      some unsigned integer value (default: 10)
    --uint64val uint64
      some unsigned integer64 value (default: 100)
    
    >> Float64-based parameters
    --floatval float64
      some float value (default: 0)
`
// Usage without application name.
const expUsageNoNameOutput = `
Usage ERROR: test error for testing usage of ` + stubApp + `

Usage:

$ ` + stubApp + ` --required-keys ... [--optional-keys ...]

    --yesno[=true|false], -y[=true|false]
      some boolean value (required option)
`

//
// Type to test AddVar() - parses date in format YYYY.MM.DD
//

type testOptVarYMD int
func (ov *testOptVarYMD) String() string {
	return fmt.Sprintf("Year %d month %d day %d", *ov / 10000, (*ov % 10000) / 100, *ov % 100)
}
func (ov *testOptVarYMD) Set(val string) error {
	ymd := make([]int, 0, 3)
	for _, part := range strings.Split(val, ".") {
		v, err := strconv.ParseInt(part, 10, 32)
		if err != nil {
			return fmt.Errorf("unparsable entry %q in date - %v", part, err)
		}
		ymd = append(ymd, int(v))
	}
	// Check length
	if len(ymd) != 3 {
		return fmt.Errorf("invalid YMD date length, want - 3, got - %d", len(ymd))
	}
	// Set decimal integer value
	*ov = testOptVarYMD(ymd[0] * 10000 + ymd[1] * 100 + ymd[2])
	// OK
	return nil
}
