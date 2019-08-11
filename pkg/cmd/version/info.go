package version

var (
	version   = "0.0.0"
	gitCommit = "$Format:%H$"          // sha1 from git, output of $(git rev-parse HEAD)
	buildDate = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)
