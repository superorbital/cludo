// The build package is intended to be used with govvv -pkg github.com/superorbital/cludo/pkg/build
package build

import "fmt"

// GitCommit is a short commit hash of the application's git source tree
var GitCommit string

// GitBranch is the current branch name the code is built off
var GitBranch string

// GitState is dirty if there are uncommitted changes, clean otherwise
var GitState string

// GitSummary is the output of git describe --tags --dirty --always
var GitSummary string

// BuildDate is a RFC3339 formatted UTC date
var BuildDate string

// Version is the version of the application
var Version string

func VersionFull() string {
	return fmt.Sprintf("%s.git-%s", Version, GitCommit)
}
