package optsparser

import (
	"time"
)

var tests = []struct{
	defaults	testOpts
	args		[]string
	want		testOpts
	needOK		bool
}{
	{
		defaults: testOpts{
			vBool:		false,
			vString:	"",
			vInt:		0,
			vInt64:		0,
			vFloat64:	0.0,
			vDuration:	time.Second,
			vUint:		0,
			vUint64:	0,
			// TODO vVar		any
		},
		args: []string{
			`--bool-opt=true`,
			`--bool-opt`,
			`--string-opt`,		`I think, therefore I am`,
			`--int-opt`,		`-430`,								// Dead Sea level below sea level
			`--int64-opt`,		`-59604644783353249`,				// Leyland prime number
			`--float64-opt`,	`3.141592`,
			`--duration-opt`,	`250560m`,							// 4176 hours => 174 days
			`--uint-opt`,		`220414`,							// The End and the Beginning
			`--uint64-opt`,		`40208000000000`,					// distance between Sun and Proxima Centauri, km
			// TODO `var-opt`, ``
		},
		want: testOpts{
			vBool:		true,
			vString:	`I think, therefore I am`,
			vInt:		-430,
			vInt64:		-59604644783353249,
			vFloat64:	3.141592,
			vDuration:	time.Hour * 4176,
			vUint:		220414,
			vUint64:	40208000000000,
			// TODO vVar		any
		},
		needOK:	true,
	},
}
