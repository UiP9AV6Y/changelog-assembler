package version

import (
	"runtime"
	"time"
)

var (
	version = "0.0.0"
	commit  = "HEAD"
	date    = "1970-01-01T00:00:00+00:00"
)

type Info struct {
	Version   string
	Commit    string
	BuildDate time.Time
	GoVersion string
	Compiler  string
	Platform  string
}

// String returns info as a human-friendly version string.
func (i Info) String() string {
	return i.Version + "-" + i.Commit
}

func NewInfo() Info {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	buildDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		buildDate = time.Unix(0, 0)
	}

	return Info{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  platform,
	}
}
