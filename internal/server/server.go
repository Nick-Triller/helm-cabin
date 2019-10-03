package server

import (
	"encoding/json"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/version"
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
	log.Printf("helm client version: %s\n", version.GetVersion())

	s.releasesChan = make(chan *rls.ListReleasesResponse)
	go PollReleases(s.releasesChan)
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
