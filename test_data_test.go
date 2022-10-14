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

// Reference value
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

// Set of tests
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

	//
	// Test FAIL cases
	//

	// Incorect option passed
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
	`11[fail]incorrect-option-bool`:		{ args: []string{`--bool-opt=invalid`}, want: testOpts{} },
	// XXX Skip testing the string type because it is impossible to pass something inappropriate for the string type
	`12[fail]incorrect-option-int`:			{ args: []string{`--int-opt`, `f430`}, want: testOpts{} },
	`13[fail]incorrect-option-int64`:		{ args: []string{`--int64-opt`, `59604644783353249.1`}, want: testOpts{} },
	`14[fail]incorrect-option-float64`:		{ args: []string{`--float64-opt`, `3.G+e0`, }, want: testOpts{} },
	`15[fail]incorrect-option-duration`:	{ args: []string{`--duration-opt`, `1Y`}, want: testOpts{} },
	`16[fail]incorrect-option-uint`:		{ args: []string{`--uint-opt`, `-220414`}, want: testOpts{} },
	`17[fail]incorrect-option-uint64`:		{ args: []string{`--uint64-opt`, `-40208000000000`}, want: testOpts{} },
	`18[fail]incorrect-option-var`:			{ args: []string{`--var-ymd-opt`, `2022/10/14`, }, want: testOpts{} },

	// Test required option missing
	`19[fail]required-opts`: {
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
}

//
// Type to test AddVar() - parses date in format YYYY.MM.DD
//

type testOptVarYMD int
func (v *testOptVarYMD) String() string {
	return fmt.Sprintf("Year %d month %d day %d", *v / 10000, (*v % 10000) / 100, *v % 100)
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
