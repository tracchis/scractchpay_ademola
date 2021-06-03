package process

import (
	"context"
	"flag"
	"fmt"
	"os"
)

// Go libraries should not set flags themselves, these are an exception.
var (
	version = flag.Bool("version", false, "display version information")
)

// Init is initialization code that provides basic functionality for processes.
//
// Init takes as parameters version information, identifying the Command Name, Semver, and a Buildstamp.
// The Buildstamp could be just a timestamp, or could include a commit hash or other reference.
//
// Init parses flags, sets up AtExit, and will start profiling if the appropriate flag is set.
//
// It returns the `context.Context` from `process.Context()`,
// and a function that `main` should `defer`,
// which will take care of executing the queued AtExit functions.
func Init(cmdname, semver, buildstamp string) (context.Context, func()) {
	buildVersion(cmdname, semver, buildstamp)

	flag.Parse()

	if *version {
		fmt.Println(Version())
		os.Exit(0)
	}

	return Context(), func() {
		runExitFuncs()

		if r := recover(); r != nil {
			panic(r)
		}
	}
}
