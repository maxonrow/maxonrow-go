package version

import "fmt"

const (
	Maj = "1"
	Min = "3"
	Fix = "9"
)

var (
	// Version is the current version of Maxonrow in string
	Version = fmt.Sprintf("%v.%v.%v", Maj, Min, Fix)

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}

func GetVersion() string {
	return Version
}
