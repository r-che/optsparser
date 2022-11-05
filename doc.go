/*
Package optsparser enhances the functionality of the standard [flag] package.

Package key features are:

 * Support long and short form of the same option
 * Supports required options to save your time from
   checking were they specified by command line or not
 * Improved Usage function - option references are displayed in the order of their addition,
   support delimiters for logical grouping of options, customizing the output of option
   specifications, displaying the flag of required options, general description and general
   form of the command to run an application

Basic usage of the package:

 func main() {
	 p := optsparser.NewParser("test-app",
		 "config-path",	// config-path is required option
	 )

	 var cfg string
	 p.AddString("config-path|c", "path to configuration", &cfg, "")
	 p.Parse()

	 fmt.Println("Configuration is:", cfg)
 }

# Option name format

The optName parameter is passed first to all methods with names beginning
with "Add" can have the following format variants:

 "long-name|l" - long option name "long-name" and short option name "l"
 "long-name"   - only long option name "long-name"
 "l"           - only short option name "l"

In case the format of optName passed to Add* function is wrong, the [OptsParser.Parse] method will panic

# Use methods of the standard flag package

Because OptsParser embeds the standard [flag] package, you can use any methods from this package.

However, you MUST NOT use the methods of the standard package that add options for parsing,
like: [flag.FlagSet.Bool], [flag.FlagSet.Float64Var], [flag.FlagSet.Var] and so on. Options added using standard package methods:

  * Do not support long/short form
  * Not printed in Usage output
  * [OptsParser.Parse] function does not consider them in parsing, what may cause error of
    missing an required option

Other functions such as [flag.FlagSet.Args], [flag.FlagSet.VisitAll],
[flag.FlagSet.Set] can be used without restrictions.

[flag]: https://pkg.go.dev/flag
[flag.FlagSet.Bool]: https://pkg.go.dev/flag#FlagSet.Bool
[flag.FlagSet.Float64Var]: https://pkg.go.dev/flag#FlagSet.Float64Var
[flag.FlagSet.Var]: https://pkg.go.dev/flag#FlagSet.Var
[flag.FlagSet.Args]: https://pkg.go.dev/flag#FlagSet.Args
[flag.FlagSet.VisitAll]: https://pkg.go.dev/flag#FlagSet.VisitAll
[flag.FlagSet.Set]: https://pkg.go.dev/flag#FlagSet.Set

# Feedback

Feel free to open the [issue] if you have any suggestions, comments or bug reports.
*/
package optsparser
