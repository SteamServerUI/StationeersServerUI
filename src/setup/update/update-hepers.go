package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

// parseVersion parses a version string (e.g., "4.6.10") into a Version struct and tries to handle a few culprits too
func parseVersion(v string) (Version, error) {
	v = strings.TrimPrefix(v, "v")
	if idx := strings.Index(v, "-"); idx != -1 {
		v = v[:idx]
	}

	var ver Version
	_, err := fmt.Sscanf(v, "%d.%d.%d", &ver.Major, &ver.Minor, &ver.Patch)
	if err != nil {
		return Version{}, fmt.Errorf("no valid X.Y.Z in tag: %s", v)
	}
	return ver, nil
}

// shouldUpdate determines if an update should proceed, returning reason if not
func shouldUpdate(current, latest Version, isInUpdateableState bool) (string, bool) {
	// Check if already up-to-date or older
	if latest.Major < current.Major ||
		(latest.Major == current.Major && latest.Minor < current.Minor) ||
		(latest.Major == current.Major && latest.Minor == current.Minor && latest.Patch <= current.Patch) {
		return "up-to-date", false
	}

	// Check if it’s a major update and not allowed
	if current.Major != latest.Major && !config.GetAllowMajorUpdates() {
		return "major-update", false
	}

	if !isInUpdateableState {
		return "not-in-updateable-state", false
	}

	return "", true
}

// getLatestRelease fetches the most recent release (or prerelease) from GitHub API
func getLatestRelease() (*githubRelease, error) {
	url := "https://api.github.com/repos/JacksonTheMaster/StationeersServerUI/releases"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response from GitHub API: %s", resp.Status)
	}

	var releases []githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub API response: %v", err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found")
	}

	// Find the most recent release
	var latestRelease *githubRelease
	var latestVersion Version
	for i, release := range releases {
		version, err := parseVersion(release.TagName)
		if err != nil {
			logger.Install.Warn(fmt.Sprintf("Skipping invalid version tag %s: %v", release.TagName, err))
			continue
		}
		if i == 0 || isReleaseNewerVersion(version, latestVersion) {
			currentVersion, err := parseVersion(config.GetVersion())
			if err == nil && isReleaseNewerVersion(currentVersion, version) {
				if release.Prerelease {
					logger.Install.Warn("Found a prerelease, but it is older than the running version. Skipping...")
					continue
				}
				continue
			}
			latestVersion = version
			latestRelease = &releases[i]
		}
	}

	if latestRelease == nil {
		return nil, fmt.Errorf("no suitable releases found")
	}

	// Log warning if the latest release is a prerelease
	if latestRelease.Prerelease && !config.GetAllowPrereleaseUpdates() {
		logger.Install.Warn(fmt.Sprintf("⚠️ Pre-release Update found: Latest version %s is a pre-release. Enable 'AllowPrereleaseUpdates' in config.json to update to it.", latestRelease.TagName))
		time.Sleep(1000 * time.Millisecond)
		logger.Install.Info("⚠️ Continuing in 3 seconds...")
		time.Sleep(1000 * time.Millisecond)
		logger.Install.Info("⚠️ Continuing in 2 seconds...")
		time.Sleep(1000 * time.Millisecond)
		logger.Install.Info("⚠️ Continuing in 1 seconds...")
		time.Sleep(1000 * time.Millisecond)
	}

	// If prerelease and AllowPrereleaseUpdates is false, find the latest stable release
	if latestRelease.Prerelease && !config.GetAllowPrereleaseUpdates() {
		var stableRelease *githubRelease
		var stableVersion Version
		for i, release := range releases {
			if release.Prerelease {
				continue
			}
			version, err := parseVersion(release.TagName)
			if err != nil {
				logger.Install.Warn(fmt.Sprintf("Skipping invalid version tag %s: %v", release.TagName, err))
				continue
			}
			if i == 0 || isReleaseNewerVersion(version, stableVersion) {
				stableVersion = version
				stableRelease = &releases[i]
			}
		}
		if stableRelease == nil {
			return nil, fmt.Errorf("no stable releases found")
		}
		return stableRelease, nil
	}

	return latestRelease, nil
}

// isNewerVersion compares two versions to determine if the first is newer
func isReleaseNewerVersion(v1, v2 Version) bool {
	if v1.Major != v2.Major {
		return v1.Major > v2.Major
	}
	if v1.Minor != v2.Minor {
		return v1.Minor > v2.Minor
	}
	if v1.Patch == v2.Patch {
		return false
	}
	return v1.Patch > v2.Patch
}
