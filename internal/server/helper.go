package server

import (
	"github.com/Nick-Triller/helm-cabin/internal/resources"
	"strconv"
)

// contains checks if a string slice contains a string value in O(n) time
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// findRelease finds a release in a ListReleasesResponse by name and version
func findRelease(releases []resources.ReleaseResource, releaseName string, version string) *resources.ReleaseResource {
	if releases == nil {
		return nil
	}
	for _, x := range releases {
		if x.Name == releaseName && strconv.Itoa(int(x.Version)) == version {
			return &x
		}
	}
	return nil
}

// findRevisions finds all revisions of a release in a ListRleasesResponse by name
func findRevisions(releases []resources.ReleaseResource, releaseName string) []resources.ReleaseResource {
	if releases == nil {
		return nil
	}
	revisions := make([]resources.ReleaseResource, 0)
	for _, x := range releases {
		if x.Name == releaseName {
			revisions = append(revisions, x)
		}
	}
	return revisions
}
