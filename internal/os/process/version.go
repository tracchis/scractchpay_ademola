package process

import (
	"os"
	"path/filepath"
)

var (
	appNameString = filepath.Base(os.Args[0])
	versionString string
)

// AppName returns the name of the application populated during process.Init().
//
// If called before process.Init() then this returns the process name reported in os.Args[0].
func AppName() string {
	return appNameString
}

// Version returns the version information populated during process.Init().
func Version() string {
	return versionString
}

func buildVersion(appName, semver, buildstamp string) {
	appNameString = appName
	versionString = appName + " " + semver

	if buildstamp != "" {
		versionString = versionString + "-" + buildstamp
	}
}
