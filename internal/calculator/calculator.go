package calculator

import (
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/Nightapes/go-semantic-release/internal/shared"

	log "github.com/sirupsen/logrus"
)

// Calculator struct
type Calculator struct{}

// New Calculator struct
func New() *Calculator {
	return &Calculator{}
}

//IncPrerelease increase prerelease by one
func (c *Calculator) IncPrerelease(preReleaseType string, version *semver.Version) (semver.Version, error) {
	defaultPrerelease := preReleaseType + ".0"
	if version.Prerelease() == "" || !strings.HasPrefix(version.Prerelease(), preReleaseType) {
		return version.SetPrerelease(defaultPrerelease)
	}

	parts := strings.Split(version.Prerelease(), ".")
	if len(parts) == 2 {
		i, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Warnf("Could not parse release tag %s, use version %s", version.Prerelease(), version.String())
			return version.SetPrerelease(defaultPrerelease)
		}
		return version.SetPrerelease(preReleaseType + "." + strconv.Itoa((i + 1)))

	}
	log.Warnf("Could not parse release tag %s, use version %s", version.Prerelease(), version.String())
	return version.SetPrerelease(defaultPrerelease)
}

//CalculateNewVersion from given commits and lastversion
func (c *Calculator) CalculateNewVersion(commits map[shared.Release][]shared.AnalyzedCommit, lastVersion *semver.Version, releaseType string, firstRelease bool) (semver.Version, bool) {
	switch releaseType {
	case "beta", "alpha":
		if len(commits["major"]) > 0 || len(commits["minor"]) > 0 || len(commits["patch"]) > 0 {
			version, _ := c.IncPrerelease(releaseType, lastVersion)
			return version, true
		}
	case "rc":
		if len(commits["major"]) > 0 || len(commits["minor"]) > 0 || len(commits["patch"]) > 0 {
			version, _ := c.IncPrerelease(releaseType, lastVersion)
			return version, false
		}
	case "release":
		if !firstRelease {
			if len(commits["major"]) > 0 {
				return lastVersion.IncMajor(), false
			} else if len(commits["minor"]) > 0 {
				return lastVersion.IncMinor(), false
			} else if len(commits["patch"]) > 0 {
				return lastVersion.IncPatch(), false
			}
		}
	}

	return *lastVersion, false
}