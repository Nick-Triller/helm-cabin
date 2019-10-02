package server

import (
	"encoding/json"
	"github.com/Nick-Triller/helmcabin/internal/releasepoller"
	"k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	releasesChan chan *rls.ListReleasesResponse
	releases *rls.ListReleasesResponse
}

// NewServer creates a server struct
func NewServer() *server {
	return &server{}
}

// Start is the main entrypoint that bootstraps the application
func (s *server) Start() {
	s.releasesChan = make(chan *rls.ListReleasesResponse)
	go releasepoller.PollReleases(s.releasesChan)
	go drainReleasesChan(s)

	http.HandleFunc("/releases", func (w http.ResponseWriter, r *http.Request) {
		filterStatuses := r.URL.Query()["status"]
		resources  := make([]*releaseResource, 0)
		for _, r := range s.releases.GetReleases() {
			// Filter according to status
			releaseStatus := strconv.FormatInt(int64(r.Info.Status.Code), 10)
			if contains(filterStatuses, releaseStatus) {
				resources = append(resources, releaseToResource(r))
			}
		}
		jsonData, _ := json.MarshalIndent(resources, "", "  ")
		_, _ = w.Write(jsonData)
	})

	addr := ":" + strconv.Itoa(8080)
	log.Println("Starting server with listen address " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func releaseToResource(r *release.Release) *releaseResource {
	resource := &releaseResource{
		Name: r.Name,
		Namespace: r.Namespace,
		Chart: &chartMetadata{
			Name:          r.Chart.Metadata.Name,
			Home:          r.Chart.Metadata.Home,
			Sources:       r.Chart.Metadata.Sources,
			Version:       r.Chart.Metadata.Version,
			Description:   r.Chart.Metadata.Description,
			Keywords:      r.Chart.Metadata.Keywords,
			Maintainers:   r.Chart.Metadata.Maintainers,
			Engine:        r.Chart.Metadata.Engine,
			Icon:          r.Chart.Metadata.Icon,
			ApiVersion:    r.Chart.Metadata.ApiVersion,
			Condition:     r.Chart.Metadata.Condition,
			Tags:          r.Chart.Metadata.Tags,
			AppVersion:    r.Chart.Metadata.ApiVersion,
			Deprecated:    r.Chart.Metadata.Deprecated,
			TillerVersion: r.Chart.Metadata.TillerVersion,
			Annotations:   r.Chart.Metadata.Annotations,
			KubeVersion:   r.Chart.Metadata.KubeVersion,
		},
		Info: &ReleaseInfo{
			Status:               &Status{
				StatusId:  r.Info.Status.Code.String(),
				Resources: r.Info.Status.Resources,
				Notes:     r.Info.Status.Notes,
			},
			FirstDeployed:        r.Info.FirstDeployed,
			LastDeployed:         r.Info.LastDeployed,
			Deleted:              r.Info.Deleted,
			Description:          r.Info.Description,
		},
		Version: r.Version,
	}
	return resource
}

func drainReleasesChan(server *server) {
	for {
		server.releases = <-server.releasesChan
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
