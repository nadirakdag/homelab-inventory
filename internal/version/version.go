package version

var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
	GoVersion = ""
)

func Set(v, c, t, g string) {
	Version = v
	Commit = c
	BuildTime = t
	GoVersion = g
}

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
}

func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
	}
}
