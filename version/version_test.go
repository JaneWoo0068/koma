package version

import (
	"regexp"
	"testing"
)

func TestLatestVersion(t *testing.T) {
	version, err := Latest()
	if err != nil {
		t.Skipf("Skipping: could not fetch latest version: %v", err)
	}

	if version == "" {
		t.Fatal("Expected a non-empty version string")
	}

	semverRegex := regexp.MustCompile(`^v?(\d+)(\.\d+){0,2}(-\w+)?$`)
	if !semverRegex.MatchString(version) {
		t.Fatalf("Expected semver format, got: %s", version)
	}
}
