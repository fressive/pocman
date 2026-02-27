package version

import "golang.org/x/mod/semver"

// Semver of Pocman Server
const SERVER_VERSION = "v1.0.0"

func MatchAgentVersion(agentVersion string) bool {
	return semver.MajorMinor(SERVER_VERSION) == semver.MajorMinor(agentVersion)
}
