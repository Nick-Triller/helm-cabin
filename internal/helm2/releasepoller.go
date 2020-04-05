package helm2

import (
	"context"
	"github.com/Nick-Triller/helm-cabin/internal/resources"
	"github.com/Nick-Triller/helm-cabin/internal/settings"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/version"
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
func PollReleases(releasesChan chan []resources.ReleaseResource, settings *settings.Settings) {
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
		} else {
			releasesChan <- convertResponseToReleaseListResources(resp)
		}

		time.Sleep(pollSleep)
	}
}

func convertResponseToReleaseListResources(resp *services.ListReleasesResponse) []resources.ReleaseResource {
	releaseResources := make([]resources.ReleaseResource, len(resp.Releases))
	for i, helm2Release := range resp.GetReleases() {
		releaseResources[i] = releaseListResourceFrom(helm2Release)
	}
	return releaseResources
}

func releaseListResourceFrom(r *release.Release) resources.ReleaseResource {
	files := make([]resources.File, len(r.Chart.Files))
	for i, helm2ChartFile := range r.Chart.Files {
		template := resources.File{
			TypeURL: helm2ChartFile.TypeUrl,
			Value:   helm2ChartFile.Value,
		}
		files[i] = template
	}

	templates := make([]resources.Template, len(r.Chart.Templates))
	for i, helm2Template := range r.Chart.Templates {
		template := resources.Template{
			Name:  helm2Template.Name,
			Data: helm2Template.Data,
		}
		templates[i] = template
	}

	maintainers := make([]resources.Maintainer, len(r.Chart.Metadata.Maintainers))
	for i, helm2Mantainer := range r.Chart.Metadata.Maintainers {
		maintainer := resources.Maintainer{
			Name:  helm2Mantainer.Name,
			Email: helm2Mantainer.Email,
			URL:   helm2Mantainer.Url,
		}
		maintainers[i] = maintainer
	}

	return resources.ReleaseResource{
		Name:      r.Name,
		Namespace: r.Namespace,
		Templates: templates,
		Files: files,
		Values: r.Chart.Values.Raw,
		Chart: &resources.ChartMetadata{
			Name:          r.Chart.Metadata.Name,
			Home:          r.Chart.Metadata.Home,
			Sources:       r.Chart.Metadata.Sources,
			Version:       r.Chart.Metadata.Version,
			Description:   r.Chart.Metadata.Description,
			Keywords:      r.Chart.Metadata.Keywords,
			Maintainers:   maintainers,
			Engine:        r.Chart.Metadata.Engine,
			Icon:          r.Chart.Metadata.Icon,
			APIVersion:    r.Chart.Metadata.ApiVersion,
			Condition:     r.Chart.Metadata.Condition,
			Tags:          r.Chart.Metadata.Tags,
			AppVersion:    r.Chart.Metadata.AppVersion,
			Deprecated:    r.Chart.Metadata.Deprecated,
			TillerVersion: r.Chart.Metadata.TillerVersion,
			Annotations:   r.Chart.Metadata.Annotations,
			KubeVersion:   r.Chart.Metadata.KubeVersion,
		},
		Info: &resources.ReleaseInfo{
			Status: &resources.Status{
				StatusID:  r.Info.Status.Code.String(),
				Notes:     r.Info.Status.Notes,
			},
			FirstDeployed: r.Info.FirstDeployed,
			LastDeployed:  r.Info.LastDeployed,
			Deleted:       r.Info.Deleted,
			Description:   r.Info.Description,
		},
		Manifest: r.Manifest,
		Version: r.Version,
	}
}
