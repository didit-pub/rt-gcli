package version

var (
	Version   = "dev"
	CommitSHA = "none"
	BuildDate = "unknown"
)

func GetVersion() string {
	return Version
}
