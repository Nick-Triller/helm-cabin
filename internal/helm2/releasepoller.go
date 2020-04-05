package helm2

import (
	"context"
	"github.com/Nick-Triller/helm-cabin/internal/settings"
	"time"

	"google.golang.org/grpc/metadata"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/version"
	log "k8s.io/klog"
)

var releaseStatuses = []release.Status_Code{
	release.Status_UNKNOWN,
	release.Status_DEPLOYED,
	release.Status_DELETED,
	release.Status_SUPERSEDED,
	release.Status_FAILED,
	release.Status_DELETING,
	release.Status_PENDING_INSTALL,
	release.Status_PENDING_UPGRADE,
	release.Status_PENDING_ROLLBACK,
}

// NewContext creates a versioned context.
func NewContext() context.Context {
	md := metadata.Pairs("x-helm-api-client", version.GetVersion())
	return metadata.NewOutgoingContext(context.Background(), md)
}

// PollReleases uses helm client to poll releasesCache
func PollReleases(releasesChan chan *services.ListReleasesResponse, settings *settings.Settings) {
	connectTiller(settings)
	pollSleep := 6 * time.Second
	for {
		listReleaseRequest := &services.ListReleasesRequest{
			Limit:  9999999999,
			Offset: "",
			// Reverse chronological
			SortBy: services.ListSort_LAST_RELEASED,
			Filter: "",
			// Reverse chronological
			SortOrder: services.ListSort_DESC,
			// Releases with any status
			StatusCodes: releaseStatuses,
			// Releases with any namespace
			Namespace: "",
		}

		resp, err := getReleasesList(listReleaseRequest)

		if err != nil {
			log.Warningf("Failed to retrieve releases: %v", err)
			return
		}

		releasesChan <- resp
		time.Sleep(pollSleep)
	}
}
