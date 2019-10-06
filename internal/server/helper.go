package server

import (
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
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
func findRelease(releases *services.ListReleasesResponse, releaseName string, version string) *release.Release {
	if releases == nil {
		return nil
	}
	for _, x := range releases.GetReleases() {
		if x.Name == releaseName {
			return x
		}
	}
	return nil
}

// findRevisions finds all revisions of a release in a ListRleasesResponse by name
func findRevisions(releases *services.ListReleasesResponse, releaseName string) []*release.Release {
	if releases == nil {
		return nil
	}
	revisions := make([]*release.Release, 0)
	for _, x := range releases.GetReleases() {
		if x.Name == releaseName {
			revisions = append(revisions, x)
		}
	}
	return revisions
}
