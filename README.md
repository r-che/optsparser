OptsParser - enhanced command line flag parser for Go
==========

[![Go Reference](https://pkg.go.dev/badge/github.com/r-che/optsparser.svg)](https://pkg.go.dev/github.com/r-che/optsparser)

The optsparser package enhances the functionality of the standard [flag] package.

-------------------------
## Installation

Install the package:

```bash
go get github.com/r-che/optsparser
```

-------------------------

## Features 

### Long and short form options

The optsparser package supports long and short form of the same option:

```go
func main() {
    p := optsparser.NewParser("test-app")

    var cfg string
    p.AddString("config-path|c", "path to configuration", &cfg, "")
    p.Parse()

    fmt.Println("Configuration is:", cfg)
}
```

Because of the features of the standard [flag] package, both forms of the option can be used with either one or two hyphens:

```
$ go run test_app.go --config-path /etc/test-app.cfg
Configuration is: /etc/test-app.cfg
$ go run test_app.go -c /etc/test-app.cfg
Configuration is: /etc/test-app.cfg
$ go run test_app.go -config-path /etc/test-app.cfg
Configuration is: /etc/test-app.cfg
$ go run test_app.go --c /etc/test-app.cfg
Configuration is: /etc/test-app.cfg
```

### Required options

The optsparser package supports required options, saving you from checking if they were specified by the command line:

```go
func main() {
    p := optsparser.NewParser("test-app",
        // List of required options
        "sleep-delay",
        "config-path",
    )

    var delay time.Duration
    p.AddDuration("sleep-delay|s", "waiting delay", &delay,
        time.Duration(-1)) // NOTE Default value, not used with the required options
    var cfg string
    p.AddString("config-path|c", "path to configuration", &cfg, "")
    var debug bool
    p.AddBool("debug|d", "enable debug output", &debug, false)
    p.Parse()  // program will exit at this point if required option is omitted

    fmt.Println("Delay", delay, "before start, configuration is:", cfg, "debug:", debug)
}
```

An example of a launch with all required options:
```
 $ go run test_app.go --sleep-delay 10s --config-path /etc/test-app.cfg 
 Delay 10s before starting, configuration is: /etc/test-app.cfg debug: false
```

An example where the required options are omitted:
```
  $ go run test_app.go -d
  Usage ERROR: required option(s) is missing: --config-path, --sleep-delay
  // Usage message is omitted
```

### Improved Usage Function

Usage Function:

  * Prints help by options in the order they were added - you can group options logically, by functional groups
  * Prints general description and general form of the command to run an application if they were specified
  * Allows to add separators between groups of options - as empty lines or as text describing the group
  * Has the ability to customize the option line output - which of the forms (short or long) is printed first,
    the separator used between them
  * Indicates whether an option is required, otherwise the default value is shown
 
### Test coverage over 99% of the code

```
$ go test -cover
PASS
coverage: 99.3% of statements
ok      github.com/r-che/optsparser     0.003s
```

[flag]: https://pkg.go.dev/flag

-------------------------

## Feedback

Feel free to open the [issue] if you have any suggestions, comments or bug reports.

[issue]: https://github.com/r-che/optsparser/issues
